package service

import (
	"context"
	"fmt"
	"strings"
)

type IUserService interface {
	// Login(request LoginRequest) (common.ResponseData, error)
	GetUser(context.Context, string) (string, error)
}

func NewUserService() IUserService {
	return new(UserService)
}

type UserService struct {
}

func (u *UserService) GetUser(_ context.Context, uid string) (string, error) {
	fmt.Println(uid)
	return strings.ToUpper(uid), nil
}
