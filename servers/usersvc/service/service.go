package service

import "github.com/Zhan9Yunhua/blog-svr/common"

type IUserService interface {
	Login(LoginRequest) (common.ResponseData, error)
	SendCode() (common.ResponseData, error)
	Register(RegisterRequest) error
	Validate(interface{}) error
	GetUser(string) (string, error)
	GetUserList() (common.ResponseData, error)
}

type UserService struct {
}

func NewUserService() *UserService {
	return new(UserService)
}
