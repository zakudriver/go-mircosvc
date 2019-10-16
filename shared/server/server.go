package server

import (
	"context"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/config"
	lg "github.com/Zhan9Yunhua/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunServer(mux *http.ServeMux, httpAddr string) {
	mux.Handle("/metrics", promhttp.Handler())

	http.Handle("/", accessControl(mux))

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	lg.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		lg.Fatal("Server Shutdown:", err)
	}
	lg.Infoln("Server exiting")

	pid := fmt.Sprintf("%d", os.Getpid())
	pidPath := config.GetConfig().PidPath

	_, openErr := os.OpenFile(pidPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr == nil {
		ioutil.WriteFile(pidPath, []byte(pid), 0)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
