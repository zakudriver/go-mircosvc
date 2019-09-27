package transport

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/config"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"time"

	userSvcEp "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoint"
	userSvcSer "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	userSvcTransport "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
)

func MakeHandler(ctx context.Context, etcdClient *etcdv3.Instancer, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	conf := config.GetConfig()
	r := mux.NewRouter()
	{
		endpoints := userSvcEp.Endponits{}
		{
			// factory, _ := userSvcFactory(ctx, "", userSvcEp.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			// endpoints.GetUserEP = factory

			factory := usersvcFactory(userSvcEp.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpointer := sd.NewEndpointer(etcdClient, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(conf.RetryMax, time.Duration(conf.RetryTimeout), balancer)
			endpoints.GetUserEP = retry
		}
		r.PathPrefix("/user").Handler(http.StripPrefix("/user", userSvcTransport.NewHTTPHandler(endpoints, tracer,
			zipkinTracer, logger)))
	}

	return r
}

func userSvcFactory(
	_ context.Context,
	addr string,
	makeEndpoint func(service userSvcSer.IUserService) endpoint.Endpoint,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) (endpoint.Endpoint, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	svc := userSvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)

	return makeEndpoint(svc), nil
}

func usersvcFactory(makeEndpoint func(service userSvcSer.IUserService) endpoint.Endpoint, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := userSvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
