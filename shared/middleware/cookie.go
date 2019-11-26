package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kum0/go-mircosvc/common"
	"github.com/kum0/go-mircosvc/shared/session"
	"net/http"
)

func CookieMiddleware(st session.Storager) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			v, ok := ctx.Value(common.SessionKey).(string)
			if !ok {
				return nil, common.NewError(http.StatusUnauthorized, "cookie 不存在.")
			}

			if se, err := st.Read(v); err == nil {
				context.WithValue(ctx, common.CookieName, se.Data)
				// if uid, ok := se.Data[common.UIDKey]; ok {
				// 	context.WithValue(ctx, common.UIDKey, uid)
				// }
				// if level, ok := se.Data[common.RoleIDKey]; ok {
				// 	context.WithValue(ctx, common.RoleIDKey, level)
				// }
				return next(ctx, request)
			}

			return nil, common.NewError(http.StatusUnauthorized, "cookie 已过期.")
		}
	}
}
