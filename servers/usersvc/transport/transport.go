package transport

import (
	"net/http"

	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kum0/blog-svr/common"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
)

func MakeHTTPHandler(eps *endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	opts := []kitTransport.ServerOption{
		kitTransport.ServerErrorEncoder(common.EncodeError),
		kitZipkin.HTTPServerTrace(zipkinTracer),
	}

	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())

	{
		handler := kitTransport.NewServer(
			eps.LoginEP,
			common.DecodeJsonRequest(new(endpoints.LoginRequest)),
			common.EncodeResponse,
			append(opts, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Login", logger)))...,
		)
		m.Handle("/login", handler).Methods("POST")
	}

	{
		handler := kitTransport.NewServer(
			eps.SendCodeEP,
			common.DecodeEmptyHttpRequest,
			common.EncodeResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "SendCode", logger)))...,
		)
		m.Handle("/code", handler).Methods("GET")
	}

	{
		handler := kitTransport.NewServer(
			eps.RegisterEP,
			common.DecodeJsonRequest(new(endpoints.RegisterRequest)),
			common.EncodeResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Register", logger)))...,
		)
		m.Handle("/register", handler).Methods("POST")
	}

	{
		handler := kitTransport.NewServer(
			eps.UserListEP,
			DecodeUserListUrlRequest,
			common.EncodeResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "UserList", logger)))...,
		)
		m.Handle("/user", handler).Methods("GET")
	}

	{
		handler := kitTransport.NewServer(
			eps.GetUserEP,
			decodeGetUserRequest,
			common.EncodeResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetUser", logger)))...,
		)
		m.Handle("/{UID}", handler).Methods("GET")
	}

	return m
}
