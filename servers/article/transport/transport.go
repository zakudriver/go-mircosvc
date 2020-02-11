package transport

import (
	"net/http"

	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kum0/go-mircosvc/common"
	"github.com/kum0/go-mircosvc/servers/article/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeHTTPHandler(
	eps *endpoints.Endpoints,
	otTracer opentracing.Tracer,
	logger log.Logger,
	opts []kitTransport.ServerOption,
) http.Handler {
	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())

	{
		handler := kitTransport.NewServer(
			eps.GetCategoriesEP,
			common.DecodeEmptyHttpRequest,
			kitTransport.EncodeJSONResponse,
			append(opts, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetCategories", logger)))...,
		)
		m.Handle("/category", handler).Methods("GET")
	}

	return m
}
