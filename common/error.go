package common

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRouteArgs           = errors.New("error route arguments")
	ErrArgs                = errors.New("error arguments")
	ErrEmpty               = errors.New("empty input")
	ErrInterfaceConversion = errors.New("interface conversion error")
)

func ArgsErr(err interface{}) error {
	switch a := err.(type) {
	case string:
		return status.Error(codes.InvalidArgument, a)
	case error:
		return status.Error(codes.InvalidArgument, a.Error())
	default:
		return status.Error(codes.InvalidArgument, "args error")
	}
}
