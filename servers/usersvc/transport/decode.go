package transport

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/kum0/blog-svr/common"
	userPb "github.com/kum0/blog-svr/pb/user"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/kum0/blog-svr/utils"
	"net/http"
)

// GerUser
func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uid, ok := vars["UID"]
	if !ok {
		return nil, common.ErrRouteArgs
	}
	return uid, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	r, ok := grpcReq.(*userPb.GetUserRequest)
	if !ok {
		return nil, errors.New("decodeGRPCGetUserRequest: interface conversion error")
	}
	return r.Uid, nil
}

func decodeGRPCGetUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	rp, ok := grpcReply.(*userPb.GetUserReply)
	if !ok {
		return nil, errors.New("decodeGRPCGetUserResponse: interface conversion error")
	}

	r := &endpoints.GetUserResponse{}
	if err := utils.StructCopy(rp, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Login
func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*userPb.LoginRequest)
	if !ok {
		return nil, errors.New("decodeGRPCLoginRequest: interface conversion error")
	}
	return &endpoints.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

func decodeGRPCLoginResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	rp, ok := grpcReply.(*userPb.LoginReply)
	if !ok {
		return nil, errors.New("decodeGRPCLoginResponse: interface conversion error")
	}

	r := &endpoints.LoginResponse{}
	if err := utils.StructCopy(rp, r); err != nil {
		return nil, err
	}
	return r, nil
}

// SendCode
func decodeGRPCSendCodeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	rp, ok := grpcReply.(*userPb.SendCodeReply)
	if !ok {
		return nil, errors.New("decodeGRPCSendCodeResponse: interface conversion error")
	}

	r := &endpoints.SendCodeResponse{}
	if err := utils.StructCopy(rp, r); err != nil {
		return nil, err
	}
	return *r, nil
}

