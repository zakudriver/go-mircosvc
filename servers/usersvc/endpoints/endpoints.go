package endpoints

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/Zhan9Yunhua/blog-svr/shared/middleware"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"

	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

type Endponits struct {
	GetUserEP endpoint.Endpoint
	// LoginEP    endpoints.Endpoint
	// SendCodeEP endpoints.Endpoint
}

func NewEndpoints(svc service.IUserService, logger log.Logger, otTracer stdopentracing.Tracer,
	zipkinTracer *zipkin.Tracer) Endponits {
	var middlewares []endpoint.Middleware
	{
		limiter := rate.NewLimiter(rate.Every(time.Second*1), 10)

		middlewares = append(
			middlewares,
			middleware.RateLimitterMiddleware(limiter),
			circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})),
		)
	}

	var getUserEndpoint endpoint.Endpoint
	{
		method := "GetUser"
		getUserEndpoint = MakeGetUserEndpoint(svc)

		mids := append(
			middlewares,
			middleware.LoggingMiddleware(log.With(logger, "method", method)),
			opentracing.TraceServer(otTracer, method),
			kitZipkin.TraceEndpoint(zipkinTracer, method),
		)

		getUserEndpoint = handleEndpointMiddleware(getUserEndpoint, mids...)
	}

	return Endponits{
		GetUserEP: getUserEndpoint,
	}
}

func (e Endponits) GetUser(ctx context.Context, uid string) (string, error) {
	r, err := e.GetUserEP(ctx, GetUserRequest{Uid: uid})
	if err != nil {
		return "", err
	}
	response := r.(GetUserRequest)
	return response.Uid, nil
}

func handleEndpointMiddleware(endpoint endpoint.Endpoint, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return endpoint
}

func MakeGetUserEndpoint(svc service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetUserRequest)
		if !ok {
			return nil, nil
		}

		name, err := svc.GetUser(ctx, req.Uid)
		if err != nil {
			return nil, err
		}

		return GetUserResponse{Name: name}, nil

	}
}
