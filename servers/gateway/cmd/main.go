package main

import (
	"context"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/transport"
	"github.com/Zhan9Yunhua/blog-svr/shared/etcd"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	"github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.GetConfig()
	log := logger.NewLogger(conf.LogPath)

	tracer := opentracing.GlobalTracer()
	zipkinTracer := zipkin.NewZipkin(log, "", "localhost:"+conf.HttpPort, conf.ServiceName)

	etcdClient := etcd.NewEtcd(conf.EtcdHost + ":" + conf.EtcdPort)

	ctx := context.Background()
	r := transport.MakeHandler(ctx, etcdClient, tracer, zipkinTracer, log)

	errs := make(chan error, 1)
	go httpServer(log, fmt.Sprintf(":%s", conf.HttpPort), r, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Info(log).Log("serviceName", conf.ServiceName, "terminated", <-errs)
}

func httpServer(lg log.Logger, addr string, handler http.Handler, errs chan error) {
	svr := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	err := svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		lg.Log("listen: %s\n", err)
	}
	errs <- err
}
