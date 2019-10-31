package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Zhan9Yunhua/blog-svr/common"
	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"net/http"
)

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

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeGRPCGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	r, ok := request.(endpoints.GetUserRequest)
	if !ok {
		return nil, errors.New("interface conversion error")
	}
	return &userPb.GetUserRequest{Uid: r.Uid}, nil
}

func encodeGRPCGetUserResponse(_ context.Context, request interface{}) (interface{}, error) {
	r, ok := request.(endpoints.GetUserResponse)
	if !ok {
		return nil, errors.New("interface conversion error")
	}
	return &userPb.GetUserReply{Uid: r.Name}, nil
}
