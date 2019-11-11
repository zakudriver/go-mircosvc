package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
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

func NewHTTPHandler(eps *endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {

	opts := []kitTransport.ServerOption{
		kitTransport.ServerErrorEncoder(encodeError),
		kitZipkin.HTTPServerTrace(zipkinTracer),
	}

	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())

	{
		handler := makeHandler(eps.LoginEP, common.DecodeCommonJsonRequest(&endpoints.LoginRequest{}),
			common.EncodeResponse,
			append(opts, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Login", logger)),
			))
		m.Handle("/login", handler).Methods("POST")
	}

	{
		handler := makeHandler(eps.SendCodeEP,
			common.DecodeEmptyHttpRequest,
			common.EncodeResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "SendCode", logger)),
			))
		m.Handle("/code", handler).Methods("GET")
	}

	{
		handler := makeHandler(eps.GetUserEP, decodeGetUserRequest, common.EncodeResponse,
			append(opts,
				kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetUser", logger)),
			))
		m.Handle("/{UID}", handler).Methods("GET")
	}

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
