package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type loginRequest struct {
	Username string `json:"username"`
}

type commonResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
	Err  string                 `json:"err,omitempty"`
}

func makeGetUserEndpoint(s UcenterServiceInterface) endpoint.Endpoint {
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
		return commonResponse{Code: 0, Msg: "ok", Data: data, Err: errmsg}, nil
	}
}
