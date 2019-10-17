package service

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/shared/middleware"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

type Endponits struct {
	GetUserEP endpoint.Endpoint
	// LoginEP    endpoints.Endpoint
	// SendCodeEP endpoints.Endpoint
}

func NewEndpoints(svc IUserService, otTracer stdopentracing.Tracer, zipkinTracer *zipkin.Tracer) Endponits {
	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = MakeGetUserEndpoint(svc)

		middlewares := make([]endpoint.Middleware, 0)
		{
			limiter := rate.NewLimiter(rate.Every(time.Second*1), 10)
			limitterMiddleware := middleware.RateLimitterMiddleware(limiter)
			middlewares = append(middlewares, limitterMiddleware)
		}

		getUserEndpoint = handleEndpointMiddleware(getUserEndpoint, middlewares...)
		// getUserEndpoint = kitZipkin.TraceEndpoint(zipkinTracer, "usersvc_GetUser")(getUserEndpoint)
		// getUserEndpoint = opentracing.TraceServer(otTracer, "usersvc_GetUser")(getUserEndpoint)
		// sumEndpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(sumEndpoint)
		// sumEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}
	return Endponits{
		GetUserEP: getUserEndpoint,
	}
}

func (e Endponits) GetUser(ctx context.Context, uid string) (string, error) {
	r, err := e.GetUserEP(ctx, uid)

	if err != nil {
		return "", err
	}
	response := r.(string)
	return response, err
}

func handleEndpointMiddleware(endpoint endpoint.Endpoint, middlewares ...endpoint.Middleware) endpoint.Endpoint {

	// endpoint = kitZipkin.TraceEndpoint(zipkinTracer, endpointName)(endpoint)

	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return endpoint
}

func MakeGetUserEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)

		name, err := svc.GetUser(ctx, req.UID)
		if err != nil {
			return nil, err
		}
		data := map[string]interface{}{
			"id": name,
		}

		return common.Response{Code: common.OK.Code(), Msg: "ok", Data: data,}, nil
	}
}
