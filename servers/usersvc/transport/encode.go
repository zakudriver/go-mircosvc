package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kum0/blog-svr/common"
	userPb "github.com/kum0/blog-svr/pb/user"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
)

func encodeResponseSetCookie(_ context.Context, w http.ResponseWriter, response interface{}) error {
	cookie := &http.Cookie{
		Name:     common.AuthHeaderKey,
		Value:    "test",
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

// Login
func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.LoginRequest)
	if !ok {
		return nil, errors.New("encodeGRPCLoginRequest: interface conversion error")
	}
	return &userPb.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

// SendCode
// ...

// Register
func encodeGRPCRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.RegisterRequest)
	if !ok {
		return nil, errors.New("encodeGRPCRegisterRequest: interface conversion error")
	}
	return &userPb.RegisterRequest{Username: req.Username, Password: req.Password, CodeID: req.CodeID}, nil
}

// UserList
func encodeGRPCUserListRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.UserListRequest)
	if !ok {
		return nil, errors.New("encodeGRPCUserListRequest: interface conversion error")
	}
	return &userPb.UserListRequest{Page: req.Page, Size: req.Size}, nil
}

func encodeGRPCUserListResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCUserListResponse: interface conversion error")
	}

	data := res.Data.(*userPb.UserListResponse)

	// us := make([]*userPb.UserResponse, 0)
	// for _, v := range data.Data {
	// 	u := new(userPb.UserResponse)
	// 	if err := utils.StructCopy(v, u); err != nil {
	// 		return nil, err
	// 	}
	// 	us = append(us, u)
	// }

	return data, nil
}
