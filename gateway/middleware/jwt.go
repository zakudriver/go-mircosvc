package middleware

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/services/session"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

func GetJwtMiddleware() endpoint.Middleware {
	secret := config.GetConfig().JwtAuthSecret

	keysFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}

	// 一个bug跟etcd冲突
	// https://github.com/coreos/etcd/issues/9357
	return kitjwt.NewParser(keysFunc, jwt.SigningMethodHS256, kitjwt.MapClaimsFactory)
}

func CookieMiddleware(st *session.Storage) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			v, ok := ctx.Value(common.SessionKey).(string)
			if !ok {
				return common.OutputResponse{
					Code: common.Error.Code(),
					Msg:  "cookie be not exists",
				}, nil
			}

			is := st.ExistsSession(v)
			if is {
				return next(ctx, request)
			}

			return common.OutputResponse{
				Code: common.Error.Code(),
				Msg:  "cookie is expired",
			}, nil
		}
	}
}
