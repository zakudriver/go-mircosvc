package zipkin

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"

	zipkinReporterHttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func NewZipkin(logger log.Logger, zipkinAddr, svcAddr, svcName string) *zipkin.Tracer {
	isNoopTracer := (zipkinAddr == "")
	reporter := zipkinReporterHttp.NewReporter(zipkinAddr)
	defer reporter.Close()

	zkEndpoint, err := zipkin.NewEndpoint(svcName, svcAddr)
	if err != nil {
		logger.Log("zipkin NewEndpoint", err)
	}
	zipkinTracer, err := zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zkEndpoint), zipkin.WithNoopTracer(isNoopTracer),
	)
	if err != nil {
		logger.Log("zipkin NewTracer", err)
		os.Exit(0)
	}
	if !isNoopTracer {
		logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinAddr)
	}

	return zipkinTracer
}
