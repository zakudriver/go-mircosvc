package service

import (
	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/shared/middleware"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func MakeHandler(us IUserService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	registerHandler := kithttp.NewServer(
		makeRegisterEndpoint(us),
		decodeRegisterRequest,
		encodeResponse,
		opts...,
	)

	authHandler := kithttp.NewServer(
		makeAuthEndpoint(us),
		decodeNoParamsRequest,
		encodeResponse,
		opts...,
	)

	userListHandler := kithttp.NewServer(
		makeGetUserListEndpoint(us),
		decodeRegisterRequest,
		encodeResponse,
		opts...,
	)

	limiter := rate.NewLimiter(rate.Every(time.Second*1), 3)
	limitterMiddleware := middleware.RateLimitterMiddleware(limiter)

	r := mux.NewRouter()
	conf := config.GetConfig()
	// 接口路由
	r.Handle(conf.Prefix+"/login", loginHandler(us, opts, limitterMiddleware)).Methods("POST")
	r.Handle(conf.Prefix+"/{UID}", getUserHandler(us, opts, limitterMiddleware)).Methods("GET")
	r.Handle(conf.Prefix+"/code", sendCodeHandler(us, opts, limitterMiddleware)).Methods("GET")
	r.Handle(conf.Prefix+"/register", registerHandler).Methods("POST")
	r.Handle(conf.Prefix+"/auth", authHandler).Methods("GET")
	r.Handle(conf.Prefix+"/list", userListHandler).Methods("GET")

	return r
}

func loginHandler(us IUserService, opts []kithttp.ServerOption, middlewares ...endpoint.Middleware) *kithttp.Server {
	endpoint := makeLoginEndpoint(us)
	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return kithttp.NewServer(
		endpoint,
		decodeLoginRequest,
		encodeResponseSetCookie,
		opts...,
	)
}

func getUserHandler(us IUserService, opts []kithttp.ServerOption, middlewares ...endpoint.Middleware) *kithttp.Server {
	endpoint := makeGetUserEndpoint(us)
	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return kithttp.NewServer(
		endpoint,
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)
}

func sendCodeHandler(us IUserService, opts []kithttp.ServerOption, middlewares ...endpoint.Middleware) *kithttp.Server {
	endpoint := makeSendCodeEndpoint(us)
	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return kithttp.NewServer(
		endpoint,
		decodeNoParamsRequest,
		encodeResponse,
		opts...,
	)
}
