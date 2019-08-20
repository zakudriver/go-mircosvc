package service

import (
	"errors"
	"fmt"
	"strings"
)

type UcenterServiceInterface interface {
	GetUser(string) (string, error)
}

type UcenterService struct{}

func (UcenterService) GetUser(s string) (string, error) {
	fmt.Printf("query:%s", s)
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

var ErrEmpty = errors.New("empty input")

type ServiceMiddleware func(UcenterServiceInterface) UcenterServiceInterface
