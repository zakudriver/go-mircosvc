package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/go-kit/kit/log"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"net/http"
)

func NewHTTPHandler(endpoints service.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	// zipkinServer := kitZipkin.HTTPServerTrace(zipkinTracer)
	//
	// options := []kitTransport.ServerOption{
	// 	kitTransport.ServerErrorEncoder(encodeError),
	// 	zipkinServer,
	// }

	opts := []kitTransport.ServerOption{
		kitTransport.ServerErrorEncoder(encodeError),
		// kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetUser", logger))
	}

	m := mux.NewRouter()
	m.Handle("/{UID}", kitTransport.NewServer(
		endpoints.GetUserEP,
		decodeGetUserRequest,
		encodeResponse,
		opts...
	)).Methods("GET")

	// m.Handle("/login", kitTransport.NewServer(
	// 	endpoints.LoginEP,
	// 	decodeLoginRequest,
	// 	encodeResponseSetCookie,
	// 	append(options, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "Login", logger)))...,
	// ))

	m.Handle("/metrics", promhttp.Handler())
	m.Handle("/test", &S{})
	return m
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(common.Response{
		Code: common.Error.Code(),
		Msg:  err.Error(),
	})
}

type S struct {
}

func (s *S) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ok")

	w.Write([]byte("ok"))
}
