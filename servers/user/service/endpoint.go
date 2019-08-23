package service

import (
	"context"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/endpoint"
)


func makeLoginEndpoint(s UserServicer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		name, err := s.GetUser(req.Username)
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

func makeGetUserEndpoint(s UcenterServiceInterface) endpoint.Endpoint {
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

		fmt.Println(req)
		return common.InnerResponse{Code: 0, Msg: "ok", Data: data, Err: errmsg}, nil
	}
}
