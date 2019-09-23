package service

import (
	"net/http"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/shared/middleware"
	sharedZipkin "github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/time/rate"
)

var (
	svrOpts = []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	zipkinTracer *zipkin.Tracer
)

func MakeHandler(us IUserService, logger kitlog.Logger) http.Handler {
	conf := config.GetConfig()

	zipkinTracer = sharedZipkin.NewZipkin(logger, conf.ZipkinAddr, conf.ServerAddr, "user_server")
	svrOpts = append(svrOpts, addZipkinSvrOpts())

	middlewares := make([]endpoint.Middleware, 0)
	{
		limiter := rate.NewLimiter(rate.Every(time.Second*1), 10)
		limitterMiddleware := middleware.RateLimitterMiddleware(limiter)
		middlewares = append(middlewares, limitterMiddleware)
	}

	r := mux.NewRouter()
	// 接口路由
	{
		// 登录
		h := makeHandler(makeLoginEndpoint(us), "login_endpoint", decodeLoginRequest,
			encodeResponseSetCookie, middlewares...)

		r.Handle(conf.Prefix+"/login", h).Methods("POST")
	}

	{
		// 根据id获取用户信息
		h := makeHandler(makeGetUserEndpoint(us), "getUser_endpoint", decodeGetUserRequest,
			encodeResponse, middlewares...)

		r.Handle(conf.Prefix+"/{UID}", h).Methods("GET")
	}

	{
		// 发送验证码
		h := makeHandler(makeSendCodeEndpoint(us), "sendCode_endpoint", decodeNoParamsRequest,
			encodeResponse, middlewares...)

		r.Handle(conf.Prefix+"/code", h).Methods("GET")
	}

	{
		// 注册
		h := makeHandler(makeRegisterEndpoint(us), "register_endpoint", decodeRegisterRequest,
			encodeResponse, middlewares...)

		r.Handle(conf.Prefix+"/register", h).Methods("POST")
	}

	{
		// 验证登录是否过期
		h := makeHandler(makeAuthEndpoint(us), "auth_endpoint", decodeNoParamsRequest,
			encodeResponse, middlewares...)

		r.Handle(conf.Prefix+"/auth", h).Methods("GET")
	}

	{
		// 获取所有用户信息
		h := makeHandler(makeGetUserListEndpoint(us), "getUserList_endpoint", decodeNoParamsRequest,
			encodeResponse, middlewares...)

		r.Handle(conf.Prefix+"/list", h).Methods("GET")
	}

	return r
}

func addZipkinSvrOpts() kithttp.ServerOption {
	if zipkinTracer != nil {
		return kitZipkin.HTTPServerTrace(zipkinTracer, kitZipkin.Name("http-transport"))
	}
	return nil
}

func makeHandler(
	endpoint endpoint.Endpoint,
	endpointName string,
	dec kithttp.DecodeRequestFunc,
	enc kithttp.EncodeResponseFunc,
	middlewares ...endpoint.Middleware) *kithttp.Server {
	endpoint = kitZipkin.TraceEndpoint(zipkinTracer, endpointName)(endpoint)

	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return kithttp.NewServer(
		endpoint,
		dec,
		enc,
		svrOpts...,
	)
}
