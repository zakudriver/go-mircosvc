package common

import (
	"context"
	"net/http"

	kitTransport "github.com/go-kit/kit/transport/http"
)

func CookieToContext() kitTransport.ServerOption {
	return kitTransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
		c, err := r.Cookie(AuthHeaderKey)
		if err != nil {
			return ctx
		}

		return context.WithValue(ctx, SessionKey, c.Value)
	})
}
