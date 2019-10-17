package service

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/middleware"
	"strings"
)

type IUserService interface {
	// Login(request LoginRequest) (common.ResponseData, error)
	GetUser(context.Context, string) (string, error)
}

func NewUserService() (s IUserService) {
	s = new(UserService)
	s = handleServiceMiddleware(s, middleware.NewPrometheusMiddleware())
	return
}

type UserService struct {
}

func (u *UserService) GetUser(_ context.Context, uid string) (string, error) {
	return strings.ToUpper(uid), nil
}

func handleServiceMiddleware(s IUserService, middlewares ...middleware.ServiceMiddleware) IUserService {
	for _, m := range middlewares {
		s = m(s)
	}

	return s
}
