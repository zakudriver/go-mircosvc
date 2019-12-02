package transport

import (
	"net/http"

	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kum0/go-mircosvc/common"
	"github.com/kum0/go-mircosvc/servers/usersvc/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeHTTPHandler(
	eps *endpoints.Endponits,
	otTracer opentracing.Tracer,
	logger log.Logger,
	opts []kitTransport.ServerOption,
) http.Handler {
	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())

	{
		handler := kitTransport.NewServer(
			eps.LoginEP,
			common.DecodeJsonRequest(new(endpoints.LoginRequest)),
			kitTransport.EncodeJSONResponse,
			append(opts, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Login", logger)))...,
		)
		m.Handle("/login", handler).Methods("POST")
	}

	{
		handler := kitTransport.NewServer(
			eps.SendCodeEP,
			common.DecodeEmptyHttpRequest,
			kitTransport.EncodeJSONResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "SendCode", logger)))...,
		)
		m.Handle("/code", handler).Methods("GET")
	}

	{
		handler := kitTransport.NewServer(
			eps.RegisterEP,
			common.DecodeJsonRequest(new(endpoints.RegisterRequest)),
			kitTransport.EncodeJSONResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Register", logger)))...,
		)
		m.Handle("/register", handler).Methods("POST")
	}

	{
		handler := kitTransport.NewServer(
			eps.UserListEP,
			DecodeUserListRequest,
			kitTransport.EncodeJSONResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "UserList", logger)))...,
		)
		m.Handle("/user", handler).Methods("GET")
	}

	{
		handler := kitTransport.NewServer(
			eps.AuthEP,
			common.DecodeEmptyHttpRequest,
			kitTransport.EncodeJSONResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Auth", logger)))...,
		)
		m.Handle("/auth", handler).Methods("GET")
	}

	{
		handler := kitTransport.NewServer(
			eps.LogoutEP,
			decodeLogoutRequest,
			kitTransport.EncodeJSONResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Logout", logger)))...,
		)
		m.Handle("/logout", handler).Methods("GET")
	}

	{
		handler := kitTransport.NewServer(
			eps.GetUserEP,
			decodeGetUserRequest,
			encodeResponseSetCookie,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetUser", logger)))...,
		)
		m.Handle("/{UID}", handler).Methods("GET")
	}

	return m
}
