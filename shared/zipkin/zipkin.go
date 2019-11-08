package zipkin

import (
	"github.com/openzipkin/zipkin-go/reporter"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinReporterHttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinReporterLog "github.com/openzipkin/zipkin-go/reporter/log"
	l "log"
)

func NewZipkin(logger log.Logger, zipkinAddr, svcAddr, svcName string) *zipkin.Tracer {
	var (
		reporter reporter.Reporter
	)

	if zipkinAddr == "" {
		reporter = zipkinReporterLog.NewReporter(l.New(os.Stderr, "", l.LstdFlags))
	} else {
		reporter = zipkinReporterHttp.NewReporter(zipkinAddr)
	}
	defer reporter.Close()

	zkEndpoint, err := zipkin.NewEndpoint(svcName, svcAddr)
	if err != nil {
		logger.Log("zipkin NewEndpoint", err)
	}
	zipkinTracer, err := zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zkEndpoint),
	)
	if err != nil {
		logger.Log("zipkin NewTracer", err)
		os.Exit(0)
	}

	return zipkinTracer
}
