package middleware

import (
	"context"
	"github.com/kum0/blog-svr/common"
	"github.com/kum0/blog-svr/shared/session"
	"github.com/go-kit/kit/endpoint"
)

func CookieMiddleware(st *session.Storage) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			v, ok := ctx.Value(common.SessionKey).(string)
			if !ok {
				return common.Response{
					Msg:  "cookie is not exists",
				}, nil
			}

			if se, err := st.ReadSession(v); err == nil {
				context.WithValue(ctx, common.AuthHeaderKey, se.Data)
				return next(ctx, request)
			}

			return common.Response{
				Msg:  "cookie is expired",
			}, nil
		}
	}
}
