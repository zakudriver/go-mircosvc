package zipkin

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"

	lg "github.com/Zhan9Yunhua/logger"
	zipkinMiddlewareHttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinReporterHttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func NewZipkin(logger log.Logger, zipkinAddr, svrAddr, svrName string) *zipkin.Tracer {
	isNoopTracer := (svrAddr == "")
	reporter := zipkinReporterHttp.NewReporter(zipkinAddr)
	defer reporter.Close()

	zEP, err := zipkin.NewEndpoint(svrName, svrAddr)
	if err != nil {
		logger.Log("zipkin NewEndpoint", err)
	}
	zipkinTracer, err := zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(isNoopTracer),
	)
	if err != nil {
		logger.Log("zipkin NewTracer",err)
	}
	if !isNoopTracer {
		logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinAddr)
	}

	return zipkinTracer
}
func NewTransport(zikkinTracer *zipkin.Tracer) http.RoundTripper {
	transport, err := zipkinMiddlewareHttp.NewTransport(zikkinTracer, zipkinMiddlewareHttp.TransportTrace(true))
	if err != nil {
		lg.Fatalln(err)
	}

	return transport
}
