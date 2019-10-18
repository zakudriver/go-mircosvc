package main

import (
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/middleware"
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
	logger := logger.NewLogger(conf.LogPath)

	tracer := opentracing.GlobalTracer()
	zipkinTracer := zipkin.NewZipkin(logger, "", "localhost:"+conf.HttpPort, conf.ServiceName)

	// etcdClient := etcd.NewEtcd(conf.EtcdHost + ":" + conf.EtcdPort)

	svc := service.NewUserService()
	svc = middleware.MakeServiceMiddleware(svc)
	ep := endpoints.NewEndpoints(svc, logger, tracer, zipkinTracer)

	handle := transport.NewHTTPHandler(ep, tracer, zipkinTracer, logger)

	errs := make(chan error, 1)
	go httpServer(logger, fmt.Sprintf(":%s", conf.HttpPort), handle, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Info(logger).Log("serviceName", conf.ServiceName, "terminated", <-errs)
}

func httpServer(logger log.Logger, addr string, handler http.Handler, errs chan error) {
	http.Handle("/", accessControl(handler))
	svr := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	err := svr.ListenAndServe()
	if err != nil {
		logger.Log("listen: %s\n", err)
	}
	errs <- err
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
