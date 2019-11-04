package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zhan9Yunhua/blog-svr/common"
	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeResponseSetCookie(_ context.Context, w http.ResponseWriter, response interface{}) error {
	cookie := &http.Cookie{
		Name:     common.AuthHeaderKey,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(common.MaxAge),
	}
	http.SetCookie(w, cookie)
	return json.NewEncoder(w).Encode(response)
}

// GetUser
func encodeGRPCGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	r, ok := request.(string)
	if !ok {
		return nil, errors.New("encodeGRPCGetUserRequest: interface conversion error")
	}

	return &userPb.GetUserRequest{Uid: r}, nil
}

func encodeGRPCGetUserResponse(_ context.Context, request interface{}) (interface{}, error) {
	r, ok := request.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCGetUserResponse: interface conversion error")
	}
	return &userPb.GetUserReply{Uid: r.Data.(string)}, nil
}

// Login
func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.LoginRequest)
	if !ok {
		return nil, errors.New("encodeGRPCLoginRequest: interface conversion error")
	}
	return &userPb.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

func encodeGRPCLoginResponse(_ context.Context, request interface{}) (interface{}, error) {
	r, ok := request.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCLoginResponse: interface conversion error")
	}
	data, ok := r.Data.(endpoints.LoginResponse)
	if !ok {
		return nil, errors.New("encodeGRPCLoginResponse: interface conversion error")
	}
	return &userPb.LoginReply{Username: data.Username}, nil
}
