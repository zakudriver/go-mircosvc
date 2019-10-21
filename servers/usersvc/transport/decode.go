package transport

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/gorilla/mux"
	"net/http"
)

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	value, ok := vars["UID"]
	if !ok {
		return nil, common.ErrRouteArgs
	}
	return endpoints.GetUserRequest{Uid: value}, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	// req := grpcReq.(*userPb.GetUserRequest)
	return endpoints.GetUserRequest{Uid: grpcReq.(string)}, nil
}
