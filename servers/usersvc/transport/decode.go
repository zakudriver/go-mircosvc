package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
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

func decodeGRPCGetUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	r, ok := grpcReply.(*userPb.GetUserReply)
	if !ok {
		return nil, errors.New("interface conversion error")
	}
	return endpoints.GetUserRequest{Uid: r.Uid}, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	r, ok := grpcReq.(*userPb.GetUserRequest)
	if !ok {
		fmt.Println("decodeGRPCGetUserRequest")
		return nil, errors.New("interface conversion error")
	}
	return endpoints.GetUserRequest{Uid: r.Uid}, nil
}
