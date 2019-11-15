package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/prometheus/client_golang/prometheus"

	userPb "github.com/kum0/blog-svr/pb/user"
)

type ServiceMiddleware func(servicer endpoints.IUserService) endpoints.IUserService

func MakeServiceMiddleware(s endpoints.IUserService) endpoints.IUserService {
	mids := []ServiceMiddleware{
		NewPrometheusMiddleware,
	}
	for _, m := range mids {
		s = m(s)
	}

	return s
}

var (
	fieldKeys = []string{"method", "error"}
)

func NewPrometheusMiddleware(next endpoints.IUserService) endpoints.IUserService {
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
	return &prometheusMiddleware{requestCount, requestLatency, countResult, next}
}

type prometheusMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	service        endpoints.IUserService
}

func (pm *prometheusMiddleware) GetUser(ctx context.Context, req string) (res *userPb.GetUserResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_user", "error", fmt.Sprint(err != nil)}
		pm.requestCount.With(lvs...).Add(1)
		pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = pm.service.GetUser(ctx, req)
	return
}

func (pm *prometheusMiddleware) timeDiff(method string, begin time.Time, err error) {
	lvs := []string{"method", method, "error", fmt.Sprint(err != nil)}
	pm.requestCount.With(lvs...).Add(1)
	pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
}

func (pm *prometheusMiddleware) Login(ctx context.Context, req endpoints.LoginRequest) (res *userPb.LoginResponse, err error) {
	defer pm.timeDiff("Login", time.Now(), err)

	res, err = pm.service.Login(ctx, req)
	return
}

func (pm *prometheusMiddleware) SendCode(ctx context.Context) (res *userPb.SendCodeResponse, err error) {
	defer pm.timeDiff("SendCode", time.Now(), err)

	res, err = pm.service.SendCode(ctx)
	return
}

func (pm *prometheusMiddleware) Register(ctx context.Context, req endpoints.RegisterRequest) (err error) {
	defer pm.timeDiff("Register", time.Now(), err)

	err = pm.service.Register(ctx, req)
	return
}

func (pm *prometheusMiddleware) UserList(ctx context.Context, req endpoints.UserListRequest) (res *userPb.
UserListResponse, err error) {
	defer pm.timeDiff("Register", time.Now(), err)

	res, err = pm.service.UserList(ctx, req)
	return
}
