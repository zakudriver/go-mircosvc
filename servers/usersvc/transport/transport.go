package transport

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/go-kit/kit/log"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"net/http"
)

func NewHTTPHandler(endpoints service.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	zipkinServer := kitZipkin.HTTPServerTrace(zipkinTracer)

	options := []kitTransport.ServerOption{
		kitTransport.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	m := mux.NewRouter()
	m.Handle("/user/{param}", kitTransport.NewServer(
		endpoints.GetUserEP,
		decodeGetUserRequest,
		encodeResponse,
		append(options, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetUser", logger)))...,
	))
	// m.Handle("/login", kitTransport.NewServer(
	// 	endpoints.LoginEP,
	// 	decodeLoginRequest,
	// 	encodeResponseSetCookie,
	// 	append(options, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Login", logger)))...,
	// ))
	return m
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(common.Response{
		Code: common.Error.Code(),
		Msg:  err.Error(),
	})
}
