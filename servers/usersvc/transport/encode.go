package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kum0/blog-svr/common"
	userPb "github.com/kum0/blog-svr/pb/user"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/kum0/blog-svr/utils"
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
	res, ok := request.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCLoginResponse: interface conversion error")
	}
	data, ok := res.Data.(endpoints.LoginResponse)
	if !ok {
		return nil, errors.New("encodeGRPCLoginResponse: interface conversion error")
	}

	r := &userPb.LoginReply{}
	if err := utils.StructCopy(data, r); err != nil {
		return nil, err
	}
	return r, nil
}

// SendCode
func encodeGRPCSendCodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	return nil, nil
}

func encodeGRPCSendCodeResponse(_ context.Context, request interface{}) (interface{}, error) {
	res, ok := request.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCSendCodeResponse: interface conversion error")
	}

	r := &userPb.SendCodeReply{}
	if err := utils.StructCopy(res, r); err != nil {
		return nil, err
	}
	return r, nil
}
