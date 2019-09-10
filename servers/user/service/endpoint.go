package service

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/endpoint"
)

type loginRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
}

func makeLoginEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (a interface{}, err error) {
		err = s.Validate(request)
		if err != nil {
			return nil, err
		}

		req := request.(loginRequest)
		name, err := s.Login(req)
		if err != nil {
			return nil, err
		}

		data := map[string]interface{}{
			"user": name,
		}

		return common.Response{Code: common.OK.Code(), Msg: "ok", Data: data,}, nil
	}
}

type registerRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
	Code     int    `json:"code" validator:"required||len=6"`
	CodeID   string `json:"codeID" validator:"required"`
}

func makeRegisterEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (a interface{}, err error) {
		err = s.Validate(request)
		if err != nil {
			return nil, err
		}

		req := request.(registerRequest)
		err = s.Register(req)
		if err != nil {
			return nil, err
		}

		return common.Response{Code: common.OK.Code(), Msg: "注册成功",}, nil
	}
}

type getUserRequest struct {
	UID string `json:"s"`
}

func makeGetUserEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
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

func makeSendCodeEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		res, err := s.SendCode()
		if err != nil {
			return nil, err
		}

		return common.Response{Msg: "注册码发送成功", Data: res,}, nil
	}
}
