package middleware

import (
	"github.com/kum0/go-mircosvc/servers/usersvc/endpoints"
)

type ServiceMiddleware func(endpoints.IUserService) endpoints.IUserService

func MakeServiceMiddleware(s endpoints.IUserService) endpoints.IUserService {
	mids := []ServiceMiddleware{
		makePrometheusMiddleware,
	}
	for _, m := range mids {
		s = m(s)
	}

	return s
}
