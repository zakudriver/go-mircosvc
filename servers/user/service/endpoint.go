package service

import (
	"context"

	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/endpoint"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func makeLoginEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		name, err := s.Login(req)

		data := map[string]interface{}{
			"user": name,
		}
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		} else {
			errmsg = ""
		}

		return common.InnerResponse{Msg: "ok", Data: data, Err: errmsg}, nil
	}
}

type getUserRequest struct {
	UID string `json:"s"`
}

func makeGetUserEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
		name, err := s.GetUser(req.UID)
		data := map[string]interface{}{
			"id": name,
		}
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		} else {
			errmsg = ""
		}

		return common.InnerResponse{ Msg: "ok", Data: data, Err: errmsg}, nil
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     int    `json:"code"`
}

func makeSendCodeEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		err := s.SendCode()
		if err != nil {
			return nil, err
		}

		return common.InnerResponse{Msg: "注册码发送成功", Data: nil, Err: nil}, nil
	}
}
