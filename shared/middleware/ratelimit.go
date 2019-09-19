package middleware

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var (
	ErrLimitExceed = errors.New("rate limit error")
)

func RateLimitterMiddleware(limitter *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limitter.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
