package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kum0/go-mircosvc/common"
	"github.com/kum0/go-mircosvc/shared/session"
)

func CookieMiddleware(st session.Storager) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			v, ok := ctx.Value(common.SessionKey).(string)
			if !ok {
				return nil, common.NewError(401, "cookie is not exists")
			}

			if se, err := st.Read(v); err == nil {
				context.WithValue(ctx, common.CookieName, se.Data)
				return next(ctx, request)
			}

			return nil, common.NewError(401, "cookie is expired")
		}
	}
}
