package endpoints

import (
	"context"
	"strings"
)

type IUserService interface {
	GetUser(context.Context, string) (string, error)
	Login(context.Context, LoginRequest) (LoginResponse, error)
}

func NewUserService() IUserService {
	return new(UserService)
}

type UserService struct {
}

func (svc *UserService) GetUser(_ context.Context, uid string) (string, error) {
	return strings.ToUpper(uid), nil
}

func (svc *UserService) Login(_ context.Context, req LoginRequest) (LoginResponse, error) {
	return LoginResponse{Username: "test"}, nil
}
