package transport

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/gorilla/mux"
	"net/http"
)

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	value, ok := vars["UID"]
	if !ok {
		return nil, common.ErrRouteArgs
	}
	return service.GetUserRequest{UID: value}, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*userPb.GetUserRequest)
	return service.GetUserRequest{UID: req.Uid}, nil
}
