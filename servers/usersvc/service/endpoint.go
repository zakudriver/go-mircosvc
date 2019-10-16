package service

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/sony/gobreaker"
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
		getUserEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserEndpoint)
		getUserEndpoint = opentracing.TraceServer(otTracer, "GetUser")(getUserEndpoint)
		getUserEndpoint = kitZipkin.TraceEndpoint(zipkinTracer, "GetUser")(getUserEndpoint)
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
