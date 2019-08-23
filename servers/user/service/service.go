package service

import (
	"errors"
	"strings"
)

type UserServicer interface {
	Login(loginRequest)
	GetUser(string) (string, error)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserService struct{}

func (UserService) GetUser(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

var ErrEmpty = errors.New("empty input")

type ServiceMiddleware func(UserServicer) UserServicer
