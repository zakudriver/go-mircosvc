package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	lg "github.com/Zhan9Yunhua/logger"
)

func RunServer(addr string, handler http.Handler) {
	svr := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	lg.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
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
