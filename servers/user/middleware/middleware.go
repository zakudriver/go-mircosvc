package middleware

import (
	"fmt"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/service"
	"github.com/go-kit/kit/metrics"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceMiddleware func(servicer service.IUserService) service.IUserService

var (
	fieldKeys = []string{"method", "error"}
)

func NewInstrumentingMiddleware() ServiceMiddleware {
	requestCount := kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "get_user",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "get_user",
		Subsystem: "user_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

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
