package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"net/http"

	userSvcEp "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoint"
	userSvcSer "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	userSvcTransport "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
)

func MakeHandler(ctx context.Context, tracer opentracing.Tracer, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	{
		endpoints := userSvcEp.Endponits{}
		{
			factory, _ := userSvcFactory(ctx, addsvc, userSvcEp.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpoints.GetUserEP = factory
		}
		r.PathPrefix("/user").Handler(http.StripPrefix("/user", userSvcTransport.NewHTTPHandler(endpoints, tracer,
			zipkinTracer, logger)))
	}

	return r
}

func userSvcFactory(
	_ context.Context,
	addsvc string,
	makeEndpoint func(service userSvcSer.IUserService) endpoint.Endpoint,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) (endpoint.Endpoint, error) {
	conn, err := grpc.Dial(addsvc, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	svc := userSvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)

	return makeEndpoint(svc), nil
}
