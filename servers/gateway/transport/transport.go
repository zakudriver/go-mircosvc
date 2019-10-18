package transport

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/servers/gateway/config"
	usersvcEndpoints "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"

	usersvcSer "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	usersvcTransport "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
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
)

func MakeHandler(ctx context.Context, etcdClient etcdv3.Client, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	conf := config.GetConfig()
	r := mux.NewRouter()
	{
		endpoints := usersvcEndpoints.Endponits{}
		ins, err := etcdv3.NewInstancer(etcdClient, "/usersvc", logger)
		if err != nil {
			logger.Log(err)
		}

		{
			factory := usersvcFactory(usersvcEndpoints.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpointer := sd.NewEndpointer(ins, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)

			retry := lb.Retry(conf.RetryMax, time.Duration(conf.RetryTimeout), balancer)
			endpoints.GetUserEP = retry
		}
		r.PathPrefix("/usersvc").Handler(http.StripPrefix("/usersvc", usersvcTransport.NewHTTPHandler(endpoints, tracer,
			zipkinTracer, logger)))
	}

	return r
}

// func userSvcFactory(
// 	_ context.Context,
// 	addr string,
// 	makeEndpoint func(service userSvcSer.IUserService) endpoints.Endpoint,
// 	tracer opentracing.Tracer,
// 	zipkinTracer *zipkin.Tracer,
// 	logger log.Logger) (endpoints.Endpoint, error) {
// 	conn, err := grpc.Dial(addr, grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	svc := userSvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)
//
// 	return makeEndpoint(svc), nil
// }

func usersvcFactory(makeEndpoint func(service usersvcSer.IUserService) endpoint.Endpoint, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := usersvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
