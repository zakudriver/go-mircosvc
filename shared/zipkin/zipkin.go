package zipkin

import (
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"os"
)

func NewZipkin(logger log.Logger, zipkinAddr, svcAddr, svcName string) (*zipkin.Tracer, reporter.Reporter) {
	var (
		isNoopTracer = (zipkinAddr == "")
		reporter     = zipkinHttp.NewReporter(zipkinAddr)
	)

	zkEndpoint, err := zipkin.NewEndpoint(svcName, svcAddr)
	if err != nil {
		logger.Log("zipkin NewEndpoint", err)
	}
	zipkinTracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(zkEndpoint),
		zipkin.WithNoopSpan(isNoopTracer),
	)
	if err != nil {
		logger.Log("zipkin NewTracer", err)
		os.Exit(0)
	}

	return zipkinTracer, reporter
}
