package service

import (
	"github.com/Zhan9Yunhua/blog-svr/common"
	"strings"
)

type UserServicer interface {
	Login(loginRequest) (string, error)
	GetUser(string) (string, error)
}

type UserService struct{}


func (UserService) GetUser(s string) (string, error) {
	if s == "" {
		return "", common.ErrEmpty
	}
	return strings.ToUpper(s), nil
}


func (UserService) Login(params loginRequest) (string, error) {
	return "", nil
}
