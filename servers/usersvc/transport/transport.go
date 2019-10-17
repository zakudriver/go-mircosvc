package transport

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
)

func NewHTTPHandler(endpoints endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {

	opts := []kitTransport.ServerOption{
		kitTransport.ServerErrorEncoder(encodeError),
		kitZipkin.HTTPServerTrace(zipkinTracer,kitZipkin.Name("usersvc-transport")),
	}

	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())
	{
		ops := append(opts,
			kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "usersvc_GetUser",
				logger)))
		m.Handle("/{UID}", makeHandler(endpoints.GetUserEP, decodeGetUserRequest, encodeResponse, ops)).Methods("GET")
	}

	// m.Handle("/login", kitTransport.NewServer(
	// 	endpoints.LoginEP,
	// 	decodeLoginRequest,
	// 	encodeResponseSetCookie,
	// 	append(options, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Login", logger)))...,
	// ))

	return m
}

func makeHandler(
	endpoint endpoint.Endpoint,
	dec kitTransport.DecodeRequestFunc,
	enc kitTransport.EncodeResponseFunc,
	ops []kitTransport.ServerOption,
) *kitTransport.Server {
	return kitTransport.NewServer(
		endpoint,
		dec,
		enc,
		ops...,
	)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(common.Response{
		Code: common.Error.Code(),
		Msg:  err.Error(),
	})
}
