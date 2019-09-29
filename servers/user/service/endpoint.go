package service

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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
	// LoginEP    endpoint.Endpoint
	// SendCodeEP endpoint.Endpoint
}

func NewEndpoints(svc IUserService, logger log.Logger, duration metrics.Histogram,
	otTracer stdopentracing.Tracer,
	zipkinTracer *zipkin.Tracer) Endponits {
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeGetUserEndpoint(svc)
		sumEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(sumEndpoint)
		sumEndpoint = opentracing.TraceServer(otTracer, "Sum")(sumEndpoint)
		sumEndpoint = kitZipkin.TraceEndpoint(zipkinTracer, "Sum")(sumEndpoint)
		// sumEndpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(sumEndpoint)
		// sumEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}
	return Endponits{
		GetUserEP: sumEndpoint,
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
