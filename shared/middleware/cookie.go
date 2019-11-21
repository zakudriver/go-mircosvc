package middleware

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/kum0/go-mircosvc/common"
	"github.com/kum0/go-mircosvc/shared/session"
)

func CookieMiddleware(st *session.Storage) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			v, ok := ctx.Value(common.SessionKey).(string)
			if !ok {
				return nil, errors.New("cookie is not exists")
			}

			if se, err := st.ReadSession(v); err == nil {
				context.WithValue(ctx, common.AuthHeaderKey, se.Data)
				return next(ctx, request)
			}

			return common.Response{
				Msg: "cookie is expired",
			}, errors.New("cookie is expired")
		}
	}
}
