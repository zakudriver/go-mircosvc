package service

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/endpoint"
)

type LoginRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
}

func makeLoginEndpoint(s IUserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (a interface{}, err error) {
		err = s.Validate(request)
		if err != nil {
			return nil, err
		}

		req := request.(LoginRequest)
		userInfo, err := s.Login(req)
		if err != nil {
			return nil, err
		}

		return common.Response{Code: common.OK.Code(), Msg: "ok", Data: userInfo,}, nil
	}
}

type RegisterRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
	Code     int    `json:"code" validator:"required||len=6"`
	CodeID   string `json:"codeID" validator:"required"`
}

func makeRegisterEndpoint(s IUserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (a interface{}, err error) {
		err = s.Validate(request)
		if err != nil {
			return nil, err
		}

		req := request.(RegisterRequest)
		err = s.Register(req)
		if err != nil {
			return nil, err
		}

		return common.Response{Code: common.OK.Code(), Msg: "注册成功",}, nil
	}
}

type GetUserRequest struct {
	UID string `json:"s"`
}

func makeGetUserEndpoint(s IUserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		name, err := s.GetUser(req.UID)
		if err != nil {
			return nil, err
		}
		data := map[string]interface{}{
			"id": name,
		}

		return common.Response{Code: common.OK.Code(), Msg: "ok", Data: data,}, nil
	}
}

func makeSendCodeEndpoint(s IUserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		res, err := s.SendCode()
		if err != nil {
			return nil, err
		}

		return common.Response{Msg: "注册码发送成功", Data: res,}, nil
	}
}

func makeAuthEndpoint(_ IUserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return common.Response{Msg: "ok",}, nil
	}
}

func makeGetUserListEndpoint(s IUserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		return common.Response{Msg: "ok",}, nil
	}
}
