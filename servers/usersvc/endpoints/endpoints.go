package endpoints

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/kum0/go-mircosvc/common"
	userPb "github.com/kum0/go-mircosvc/pb/user"
	"github.com/kum0/go-mircosvc/shared/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Endponits struct {
	GetUserEP  endpoint.Endpoint
	LoginEP    endpoint.Endpoint
	SendCodeEP endpoint.Endpoint
	RegisterEP endpoint.Endpoint
	UserListEP endpoint.Endpoint
	AuthEP     endpoint.Endpoint
	LogoutEP   endpoint.Endpoint
}

func (e *Endponits) GetUser(ctx context.Context, uid string) (res *userPb.GetUserResponse, err error) {
	r, err := e.GetUserEP(ctx, uid)
	if r != nil {
		res = r.(*userPb.GetUserResponse)
	}
	return
}

func (e *Endponits) Login(ctx context.Context, request LoginRequest) (res *userPb.LoginResponse, err error) {
	r, err := e.LoginEP(ctx, request)
	if r != nil {
		res = r.(*userPb.LoginResponse)
	}
	return
}

func (e *Endponits) SendCode(ctx context.Context) (res *userPb.SendCodeResponse, err error) {
	r, err := e.SendCodeEP(ctx, nil)
	if r != nil {
		res = r.(*userPb.SendCodeResponse)
	}
	return
}

func (e *Endponits) Register(ctx context.Context, request RegisterRequest) error {
	_, err := e.RegisterEP(ctx, request)
	return err
}

func (e *Endponits) UserList(ctx context.Context, request UserListRequest) (res *userPb.UserListResponse, err error) {
	r, err := e.UserListEP(ctx, request)
	if r != nil {
		res = r.(*userPb.UserListResponse)
	}
	return
}

func (e *Endponits) Logout(ctx context.Context, req LogoutRequest) error {
	_, err := e.LogoutEP(ctx, req)
	return err
}

func NewEndpoints(svc IUserService, logger log.Logger, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer) *Endponits {

	return &Endponits{
		GetUserEP:  makeEndpoint(MakeGetUserEndpoint(svc), "GetUser", logger, otTracer, zipkinTracer),
		LoginEP:    makeEndpoint(MakeLoginEndpoint(svc), "Login", logger, otTracer, zipkinTracer),
		SendCodeEP: makeEndpoint(MakeSendCodeEndpoint(svc), "SendCode", logger, otTracer, zipkinTracer),
		RegisterEP: makeEndpoint(MakeRegisterEndpoint(svc), "Register", logger, otTracer, zipkinTracer),
		UserListEP: makeEndpoint(MakeUserListEndpoint(svc), "UseList", logger, otTracer, zipkinTracer),
		LogoutEP:   makeEndpoint(MakeLogoutEndpoint(svc), "Logout", logger, otTracer, zipkinTracer),
	}
}

// GetUser
func MakeGetUserEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(string)
		if !ok {
			return nil, errors.New("MakeGetUserEndpoint: interface conversion error")
		}

		res, err := svc.GetUser(ctx, req)

		return common.Response{Data: res}, err
	}
}

// Login
func MakeLoginEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*LoginRequest)
		if !ok {
			return nil, errors.New("MakeLoginEndpoint: interface conversion error")
		}

		res, err := svc.Login(ctx, *req)

		header := make(map[string][]string)
		if res != nil {
			header["Set-Cookie"] = []string{res.Cookie}
		}

		return common.Response{Data: res, Msg: "登陆成功.", Header: header}, err
	}
}

// SendCode
func MakeSendCodeEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := svc.SendCode(ctx)

		return common.Response{Data: res, Msg: "验证码发送成功."}, err
	}
}

// Register
func MakeRegisterEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*RegisterRequest)
		if !ok {
			return nil, errors.New("MakeRegisterEndpoint: interface conversion error")
		}

		err := svc.Register(ctx, *req)

		return common.Response{Msg: "注册成功"}, err
	}
}

// UserList
func MakeUserListEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*UserListRequest)
		if !ok {
			return nil, errors.New("MakeUserListEndpoint: interface conversion error")
		}

		res, err := svc.UserList(ctx, *req)

		return common.Response{Data: res}, err
	}
}

// Auth
func MakeAuthEndpoint(_ IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return common.Response{Msg: "auth ok"}, nil
	}
}

// Logout
func MakeLogoutEndpoint(svc IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*LogoutRequest)
		if !ok {
			return nil, errors.New("MakeLogoutEndpoint: interface conversion error")
		}

		err := svc.Logout(ctx, *req)
		return common.Response{Msg: "logout ok"}, err
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
