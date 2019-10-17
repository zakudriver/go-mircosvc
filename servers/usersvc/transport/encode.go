package transport

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"net/http"

	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
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

func encodeGRPCGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.GetUserRequest)
	return &userPb.GetUserRequest{Uid: resp.UID}, nil
}
