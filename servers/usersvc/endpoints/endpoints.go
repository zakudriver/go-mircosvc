package endpoints

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/kum0/blog-svr/common"
	"github.com/kum0/blog-svr/shared/middleware"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

type Endponits struct {
	GetUserEP  endpoint.Endpoint
	LoginEP    endpoint.Endpoint
	SendCodeEP endpoint.Endpoint
}

func (e *Endponits) GetUser(ctx context.Context, uid string) (string, error) {
	r, err := e.GetUserEP(ctx, uid)
	if err != nil {
		return "", err
	}

	return r.(string), nil
}

func (e *Endponits) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	r, err := e.LoginEP(ctx, request)
	if err != nil {
		return LoginResponse{}, err
	}
	return r.(LoginResponse), nil
}

func (e *Endponits) SendCode(ctx context.Context) error {
	_, err := e.SendCodeEP(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewEndpoints(svc IUserService, logger log.Logger, otTracer stdopentracing.Tracer,
	zipkinTracer *zipkin.Tracer) *Endponits {
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

	var loginEndpoint endpoint.Endpoint
	{
		method := "Login"
		loginEndpoint = MakeLoginEndpoint(svc)

		mids := append(
			middlewares,
			middleware.LoggingMiddleware(log.With(logger, "method", method)),
			opentracing.TraceServer(otTracer, method),
			kitZipkin.TraceEndpoint(zipkinTracer, method),
		)
		loginEndpoint = handleEndpointMiddleware(loginEndpoint, mids...)
	}

	return &Endponits{
		GetUserEP: getUserEndpoint,
		LoginEP:   loginEndpoint,
	}
}

func MakeGetUserEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetUserRequest)
		if !ok {
			return nil, errors.New("MakeGetUserEndpoint: interface conversion error")
		}

		name, err := svc.GetUser(ctx, req.Uid)
		if err != nil {
			return nil, err
		}

		return common.Response{Data: name}, nil
	}
}

func MakeLoginEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(LoginRequest)
		if !ok {
			return nil, errors.New("MakeLoginEndpoint: interface conversion error")
		}

		res, err := svc.Login(ctx, req)
		if err != nil {
			return nil, errors.New("login error")
		}

		return common.Response{Data: res}, nil
	}
}

func MakeSendCodeEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(LoginRequest)
		if !ok {
			return nil, errors.New("MakeLoginEndpoint: interface conversion error")
		}

		res, err := svc.Login(ctx, req)
		if err != nil {
			return nil, errors.New("login error")
		}

		return common.Response{Data: res}, nil
	}
}

func handleEndpointMiddleware(endpoint endpoint.Endpoint, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return endpoint
}
