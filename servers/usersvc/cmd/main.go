package main

import (
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	"github.com/opentracing/opentracing-go"
)

func main() {
	conf := config.GetConfig()
	log := logger.NewLogger(conf.LogPath)

	tracer := opentracing.GlobalTracer()
	zipkinTracer := zipkin.NewZipkin(log, conf.ZipkinAddr, "localhost:"+conf.HttpPort, conf.ServiceName)

	// etcdClient := etcd.NewEtcd(conf.EtcdHost + ":" + conf.EtcdPort)

	svc := service.NewUserService()
	ep := service.NewEndpoints(svc, tracer, zipkinTracer)

	handle := transport.NewHTTPHandler(ep, tracer, zipkinTracer, log)

	errs := make(chan error, 1)
	go httpServer(log, fmt.Sprintf(":%s", conf.HttpPort), handle, errs)

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
