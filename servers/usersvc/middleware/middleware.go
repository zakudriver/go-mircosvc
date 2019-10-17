package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/go-kit/kit/metrics"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceMiddleware func(servicer service.IUserService) service.IUserService

func MakeServiceMiddleware(s service.IUserService) service.IUserService {
	mids := []ServiceMiddleware{
		NewPrometheusMiddleware(),
	}
	for _, m := range mids {
		s = m(s)
	}

	return s
}

var (
	fieldKeys = []string{"method", "error"}
)

func NewPrometheusMiddleware() ServiceMiddleware {
	requestCount := kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "user_space",
		Subsystem: "usersvc",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "user_space",
		Subsystem: "usersvc",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "user_space",
		Subsystem: "usersvc",
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

func (pm prometheusMiddleware) GetUser(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_user", "error", fmt.Sprint(err != nil)}
		pm.requestCount.With(lvs...).Add(1)
		pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = pm.IUserService.GetUser(ctx, s)
	return
}

// func (pm prometheusMiddleware) Login(params endpoints.LoginRequest) (data common.ResponseData, err error) {
// 	defer func(begin time.Time) {
// 		lvs := []string{"method", "login", "error", fmt.Sprint(err != nil)}
// 		pm.requestCount.With(lvs...).Add(1)
// 		pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
// 	}(time.Now())
//
// 	data, err = pm.IUserService.Login(params)
// 	return
// }
