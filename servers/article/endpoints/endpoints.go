package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/kum0/go-mircosvc/common"
	"github.com/kum0/go-mircosvc/shared/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
)

type Endpoints struct {
	GetCategoriesEP endpoint.Endpoint
}

func (e *Endpoints) GetCategories(ctx context.Context) (error, error) {
	r, err := e.GetCategoriesEP(ctx, nil)

	if r != nil {
		return nil, nil
	}

	return err, nil
}

func NewEndpoints(svc ArticleServicer, logger log.Logger, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer) *Endpoints {

	return &Endpoints{
		GetCategoriesEP: makeEndpoint(MakeGetCategoriesEndpoint(svc), "MakeGetCategories", logger, otTracer, zipkinTracer),
	}
}

func MakeGetCategoriesEndpoint(svc ArticleServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := svc.GetCategories(ctx)

		return common.Response{Data: res}, err
	}
}

func makeEndpoint(
	ep endpoint.Endpoint,
	method string,
	logger log.Logger,
	otTracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	middlewares ...endpoint.Middleware,
) endpoint.Endpoint {
	limiter := rate.NewLimiter(rate.Every(time.Second*1), 10)

	middlewares = append(
		middlewares,
		middleware.RateLimitterMiddleware(limiter),
		circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})),
		kitOpentracing.TraceServer(otTracer, method),
		kitZipkin.TraceEndpoint(zipkinTracer, method),
		middleware.LoggingMiddleware(log.With(logger, "method", method)),
	)

	for _, m := range middlewares {
		ep = m(ep)
	}

	return ep
}
