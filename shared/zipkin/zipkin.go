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
	reporter := zipkinReporterHttp.NewReporter(svrAddr)
	defer reporter.Close()

	zEP, _ := zipkin.NewEndpoint(svrName, zipkinAddr)
	zipkinTracer, err := zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(isNoopTracer),
	)
	if err != nil {
		logger.Log(err)
	}
	if !isNoopTracer {
		logger.Log("tracer", "Zipkin", "type", "Native", "URL", svrAddr)
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

