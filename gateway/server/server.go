package server

import (
	"context"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/gateway/router"
	lg "github.com/Zhan9Yunhua/logger"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunServer(addr string, router *router.Router) {
	srv := &http.Server{
		Addr:    addr,
		Handler: router.R,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatalf("listen: %s\n", err)
		}
	}()

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
