package middleware

import (
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/go-kit/kit/metrics"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceMiddleware func(servicer service.IUserService) service.IUserService

var (
	fieldKeys = []string{"method", "error"}
)

func NewPrometheusMiddleware() ServiceMiddleware {
	requestCount := kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "user_space",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "user_space",
		Subsystem: "user_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "user_space",
		Subsystem: "user_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	return func(next service.IUserService) service.IUserService {
		return prometheusMiddleware{requestCount, requestLatency, countResult, next}
	}
}

type prometheusMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	service.IUserService
}

func (pm prometheusMiddleware) GetUser(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_user", "error", fmt.Sprint(err != nil)}
		pm.requestCount.With(lvs...).Add(1)
		pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = pm.IUserService.GetUser(s)
	return
}

func (pm prometheusMiddleware) Login(params service.LoginRequest) (data common.ResponseData, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login", "error", fmt.Sprint(err != nil)}
		pm.requestCount.With(lvs...).Add(1)
		pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, err = pm.IUserService.Login(params)
	return
}
