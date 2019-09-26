package main

import (
	"context"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/transport"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	"github.com/opentracing/opentracing-go"
)

func main() {
	log := logger.NewLogger("")

	tracer := opentracing.GlobalTracer()
	zipkinTracer := zipkin.NewZipkin(log, "", "", "")

	ctx := context.Background()
	r := transport.MakeHandler(ctx, tracer, zipkinTracer, log)
	runServer(log, "", r)
}

func runServer(lg log.Logger, addr string, handler http.Handler) {
	svr := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Log("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	lg.Log("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		lg.Log("Server Shutdown:", err)
	}
	lg.Log("Server exiting")

	pid := fmt.Sprintf("%d", os.Getpid())

	_, openErr := os.OpenFile("", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr == nil {
		ioutil.WriteFile("", []byte(pid), 0)
	}
}
