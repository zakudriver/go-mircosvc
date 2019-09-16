package middleware

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"time"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/service"
)

type ServiceMiddleware func(servicer service.IUserService) service.IUserService

var (
	fieldKeys    = []string{"method", "error"}
	requestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here
)

func InstrumentingMiddleware() ServiceMiddleware {
	return func(next service.IUserService) service.IUserService {
		return instrumentingMiddleware{requestCount, requestLatency, countResult, next}
	}
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	service.IUserService
}

func (mw instrumentingMiddleware) GetUser(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.IUserService.GetUser(s)
	return
}
