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

		return common.InnerResponse{Code: 0, Msg: "ok", Data: data, Err: errmsg}, nil
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

		return common.InnerResponse{Code: 0, Msg: "ok", Data: data, Err: errmsg}, nil
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     int    `json:"code"`
}
