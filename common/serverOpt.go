package common

import (
	"context"
	kitTransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

func CookieToContext() kitTransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		c, err := r.Cookie(AuthHeaderKey)
		if err != nil {
			return ctx
		}

		return context.WithValue(ctx, SessionKey, c.Value)
	}
}
