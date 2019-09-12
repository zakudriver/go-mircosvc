package service

import (
	"net/http"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(bs IUserService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	loginHandler := kithttp.NewServer(
		makeLoginEndpoint(bs),
		decodeLoginRequest,
		encodeResponseSetCookie,
		opts...,
	)

	getUserHandler := kithttp.NewServer(
		makeGetUserEndpoint(bs),
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)

	sendCodeHandler := kithttp.NewServer(
		makeSendCodeEndpoint(bs),
		decodeNoParamsRequest,
		encodeResponse,
		opts...,
	)

	registerHandler := kithttp.NewServer(
		makeRegisterEndpoint(bs),
		decodeRegisterRequest,
		encodeResponse,
		opts...,
	)

	authHandler := kithttp.NewServer(
		makeAuthEndpoint(bs),
		decodeNoParamsRequest,
		encodeResponse,
		opts...,
	)

	userListHandler := kithttp.NewServer(
		makeGetUserListEndpoint(bs),
		decodeRegisterRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	conf := config.GetConfig()
	// 接口路由
	r.Handle(conf.Prefix+"/login", loginHandler).Methods("POST")
	r.Handle(conf.Prefix+"/{UID}", getUserHandler).Methods("GET")
	r.Handle(conf.Prefix+"/code", sendCodeHandler).Methods("GET")
	r.Handle(conf.Prefix+"/register", registerHandler).Methods("POST")
	r.Handle(conf.Prefix+"/auth", authHandler).Methods("GET")
	r.Handle(conf.Prefix+"/list", userListHandler).Methods("GET")

	return r
}
