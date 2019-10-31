package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/transport"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	sharedZipkin "github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	sharedEtcd "github.com/Zhan9Yunhua/blog-svr/shared/etcd"
	kitTransportGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/openzipkin/zipkin-go"
	zipkinGrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
)

const grpcRouterReg = `([a-zA-Z]+)/`

func main() {
	conf := config.GetConfig()
	log := logger.NewLogger(conf.LogPath)

	tracer := opentracing.GlobalTracer()
	zipkinTracer := sharedZipkin.NewZipkin(log, "", "localhost:"+conf.HttpPort, conf.ServiceName)
	etcdClient := sharedEtcd.NewEtcd(conf.EtcdAddr)

	ctx := context.Background()
	r := transport.MakeHandler(ctx, etcdClient, tracer, zipkinTracer, log)

	errs := make(chan error, 2)
	go httpServer(log, conf.HttpPort, r, errs)
	go grpcServer(zipkinTracer, conf.GrpcPort, log, errs)

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

func grpcServer(zipkinTracer *zipkin.Tracer, port string, logger log.Logger, errs chan error) {
	if port == "" {
		return
	}
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		level.Error(logger).Log("GRPC", "proxy", "listen", port, "err", err)
		os.Exit(1)
	}

	routerMap := map[string]string{
		"usersvc": "localhost:5002",
	}

	re := regexp.MustCompile(grpcRouterReg)
	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		serviceName := func(fullMethodName string) string {
			x := re.FindSubmatch([]byte(fullMethodName))
			return strings.ToLower(string(x[1]))
		}(fullMethodName)

		// Make sure we never forward internal services.
		if _, ok := routerMap[serviceName]; !ok {
			return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		// Copy the inbound metadata explicitly.
		outCtx, _ := context.WithCancel(ctx)
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())

		if ok {
			conn, err := grpc.DialContext(
				ctx,
				routerMap[serviceName],
				grpc.WithInsecure(),
				grpc.WithStatsHandler(zipkinGrpc.NewClientHandler(zipkinTracer)),
				grpc.WithDefaultCallOptions(grpc.CallCustomCodec(proxy.Codec()), grpc.FailFast(false)),
			)
			return outCtx, conn, err
		}
		return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
	}

	var server *grpc.Server
	level.Info(logger).Log("GRPC", "proxy", "exposed", port)
	server = grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
		grpc.UnaryInterceptor(kitTransportGrpc.Interceptor),
		grpc.StatsHandler(zipkinGrpc.NewServerHandler(zipkinTracer)),
	)
	reflection.Register(server)
	errs <- server.Serve(listener)
}
