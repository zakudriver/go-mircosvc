package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kum0/blog-svr/servers/gateway/config"
	"github.com/kum0/blog-svr/servers/gateway/transport"
	sharedEtcd "github.com/kum0/blog-svr/shared/etcd"
	"github.com/kum0/blog-svr/shared/logger"
	sharedZipkin "github.com/kum0/blog-svr/shared/zipkin"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
)

func main() {
	conf := config.GetConfig()
	log := logger.NewLogger(conf.LogPath)

	zipkinTracer, reporter := sharedZipkin.NewZipkin(log, conf.ZipkinAddr, "localhost:"+conf.HttpPort,
		conf.ServiceName)
	defer reporter.Close()

	opentracing.SetGlobalTracer(zipkinot.Wrap(zipkinTracer))
	tracer := opentracing.GlobalTracer()

	etcdClient := sharedEtcd.NewEtcd(conf.EtcdAddr)

	r := transport.MakeHandler(etcdClient, tracer, zipkinTracer, log, conf.RetryMax, conf.RetryTimeout)

	errs := make(chan error, 1)
	go httpServer(log, conf.HttpPort, r, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Info(log).Log("serviceName", conf.ServiceName, "terminated", <-errs)
}

func httpServer(lg log.Logger, port string, handler http.Handler, errs chan error) {
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
	err := svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		lg.Log("listen: %s\n", err)
	}
	errs <- err
}
