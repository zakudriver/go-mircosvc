package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kum0/blog-svr/shared/db"
	"github.com/kum0/blog-svr/shared/email"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kum0/blog-svr/servers/usersvc/config"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/kum0/blog-svr/servers/usersvc/middleware"
	"github.com/kum0/blog-svr/servers/usersvc/transport"
	sharedEtcd "github.com/kum0/blog-svr/shared/etcd"
	"github.com/kum0/blog-svr/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	kitGrpc "github.com/go-kit/kit/transport/grpc"
	userPb "github.com/kum0/blog-svr/pb/user"
	sharedZipkin "github.com/kum0/blog-svr/shared/zipkin"
	"github.com/opentracing/opentracing-go"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	conf := config.GetConfig()
	logger := logger.NewLogger(conf.LogPath)

	tracer := opentracing.GlobalTracer()
	zipkinTracer := sharedZipkin.NewZipkin(logger, "", "localhost:"+conf.HttpPort, conf.ServiceName)

	{
		etcdClient := sharedEtcd.NewEtcd(conf.EtcdAddr)
		register := sharedEtcd.Register("/usersvc", "localhost:"+conf.GrpcPort, etcdClient, logger)
		defer register.Register()
	}

	var svc endpoints.IUserService
	{
		mdb := db.NewMysql("", "", "", "")
		rd := db.NewRedis("", "", 0, 0)
		email := email.NewEmail("", "", "", "", 0)
		svc = endpoints.NewUserService(mdb, rd, email)
		svc = middleware.MakeServiceMiddleware(svc)
	}
	ep := endpoints.NewEndpoints(svc, logger, tracer, zipkinTracer)

	hs := health.NewServer()
	hs.SetServingStatus(conf.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	errs := make(chan error, 1)
	go grpcServer(transport.MakeGRPCServer(ep, tracer, zipkinTracer, logger), conf.GrpcPort, hs, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Info(logger).Log("serviceName", conf.ServiceName, "terminated", <-errs)
}

func grpcServer(grpcsvc userPb.UsersvcServer, port string, hs *health.Server, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		level.Error(logger).Log("protocol", "GRPC", "listen", port, "err", err)
		os.Exit(1)
	}
	level.Info(logger).Log("protocol", "GRPC", "protocol", "GRPC", "exposed", port)

	server := grpc.NewServer(grpc.UnaryInterceptor(kitGrpc.Interceptor))
	userPb.RegisterUsersvcServer(server, grpcsvc)
	healthgrpc.RegisterHealthServer(server, hs)
	reflection.Register(server)
	errs <- server.Serve(listener)
}

func httpServer(logger log.Logger, port string, handler http.Handler, errs chan error) {
	http.Handle("/", accessControl(handler))

	p := fmt.Sprintf(":%s", port)
	svr := &http.Server{
		Addr:    p,
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
