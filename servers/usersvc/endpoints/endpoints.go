package endpoints

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/kum0/blog-svr/common"
	"github.com/kum0/blog-svr/shared/middleware"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Endponits struct {
	GetUserEP  endpoint.Endpoint
	LoginEP    endpoint.Endpoint
	SendCodeEP endpoint.Endpoint
	RegisterEP endpoint.Endpoint
}

func (e *Endponits) GetUser(ctx context.Context, uid string) (*GetUserResponse, error) {
	res, err := e.GetUserEP(ctx, uid)
	if err != nil {
		return nil, err
	}

	return res.(*GetUserResponse), nil
}

func (e *Endponits) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	res, err := e.LoginEP(ctx, request)
	if err != nil {
		return nil, err
	}
	return res.(*LoginResponse), nil
}

func (e *Endponits) SendCode(ctx context.Context) (*SendCodeResponse, error) {
	res, err := e.SendCodeEP(ctx, nil)
	if err != nil {
		return nil, err
	}
	return res.(*SendCodeResponse), nil
}

func (e *Endponits) Register(ctx context.Context, request RegisterRequest) error {
	_, err := e.RegisterEP(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func NewEndpoints(svc IUserService, logger log.Logger, otTracer stdopentracing.Tracer, zipkinTracer *zipkin.Tracer) *Endponits {
	var middlewares []endpoint.Middleware
	{
		limiter := rate.NewLimiter(rate.Every(time.Second*1), 10)

		middlewares = append(
			middlewares,
			middleware.RateLimitterMiddleware(limiter),
			circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})),
		)
	}

	return &Endponits{
		GetUserEP:  makeEndpoint(MakeGetUserEndpoint(svc), "GetUser", logger, otTracer, zipkinTracer, middlewares),
		LoginEP:    makeEndpoint(MakeLoginEndpoint(svc), "Login", logger, otTracer, zipkinTracer, middlewares),
		SendCodeEP: makeEndpoint(MakeSendCodeEndpoint(svc), "SendCode", logger, otTracer, zipkinTracer, middlewares),
		RegisterEP: makeEndpoint(MakeRegisterEndpoint(svc), "Register", logger, otTracer, zipkinTracer, middlewares),
	}
}

func MakeGetUserEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(string)
		if !ok {
			return nil, errors.New("MakeGetUserEndpoint: interface conversion error")
		}

		r, err := svc.GetUser(ctx, req)
		if err != nil {
			return nil, err
		}

		return common.Response{Data: r}, nil
	}
}

func MakeLoginEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*LoginRequest)
		if !ok {
			return nil, errors.New("MakeLoginEndpoint: interface conversion error")
		}

		res, err := svc.Login(ctx, *req)
		if err != nil {
			return nil, err
		}

		return common.Response{Data: res}, nil
	}
}

func MakeSendCodeEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := svc.SendCode(ctx)
		if err != nil {
			return nil, err
		}

		return common.Response{Data: res, Msg: "验证码发送成功"}, nil
	}
}

func MakeRegisterEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*RegisterRequest)
		if !ok {
			return nil, errors.New("MakeRegisterEndpoint: interface conversion error")
		}

		err := svc.Register(ctx, *req)
		if err != nil {
			return nil, err
		}

		return common.Response{Msg: "注册成功"}, nil
	}
}

func makeEndpoint(
	ep endpoint.Endpoint,
	method string,
	logger log.Logger,
	otTracer stdopentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	middlewares []endpoint.Middleware,
) endpoint.Endpoint {

	mids := append(
		middlewares,
		opentracing.TraceServer(otTracer, method),
		kitZipkin.TraceEndpoint(zipkinTracer, method),
		middleware.LoggingMiddleware(log.With(logger, "method", method)),
	)

	for _, m := range mids {
		ep = m(ep)
	}

	return ep
}
