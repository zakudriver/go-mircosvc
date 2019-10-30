package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/middleware"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
	sharedEtcd "github.com/Zhan9Yunhua/blog-svr/shared/etcd"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	sharedZipkin "github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	kitGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	conf := config.GetConfig()
	logger := logger.NewLogger(conf.LogPath)

	tracer := opentracing.GlobalTracer()
	zipkinTracer := sharedZipkin.NewZipkin(logger, "", "localhost:"+conf.HttpPort, conf.ServiceName)

	etcdClient := sharedEtcd.NewEtcd(conf.EtcdHost + ":" + conf.EtcdPort)
	httpRegister := sharedEtcd.Register("/usersvc", "localhost:5001", etcdClient, logger)
	grpcRegister := sharedEtcd.Register("/usersvc", "localhost:5002", etcdClient, logger)
	defer httpRegister.Register()
	defer grpcRegister.Register()

	svc := service.NewUserService()
	svc = middleware.MakeServiceMiddleware(svc)
	ep := endpoints.NewEndpoints(svc, logger, tracer, zipkinTracer)

	handle := transport.NewHTTPHandler(ep, tracer, zipkinTracer, logger)

	errs := make(chan error, 2)

	hs := health.NewServer()
	hs.SetServingStatus(conf.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	go httpServer(logger, fmt.Sprintf(":%s", conf.HttpPort), handle, errs)
	go grpcServer(ep, tracer, zipkinTracer, conf.GrpcPort, hs, logger, errs)

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

func grpcServer(endpoints endpoints.Endponits, tracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	port string, hs *health.Server, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		level.Error(logger).Log("protocol", "GRPC", "listen", port, "err", err)
		os.Exit(1)
	}

	var server *grpc.Server
	level.Info(logger).Log("protocol", "GRPC", "protocol", "GRPC", "exposed", port)
	server = grpc.NewServer(grpc.UnaryInterceptor(kitGrpc.Interceptor))
	userPb.RegisterUsersvcServer(server, transport.MakeGRPCServer(endpoints, tracer, zipkinTracer, logger))
	healthgrpc.RegisterHealthServer(server, hs)
	reflection.Register(server)
	errs <- server.Serve(listener)
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
