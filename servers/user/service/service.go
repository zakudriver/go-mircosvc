package service

import (
	"context"
	"strings"
)

type IUserService interface {
	// Login(request LoginRequest) (common.ResponseData, error)
	GetUser(context.Context, string) (string, error)
}

type UserService struct {
}

func NewUserService() *UserService {
	return new(UserService)
}

func GetUser(_ context.Context, uid string) (string, error) {
	return strings.ToUpper(uid), nil
}
