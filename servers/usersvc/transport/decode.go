package transport

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	value, ok := vars["UID"]
	if !ok {
		return nil, common.ErrRouteArgs
	}
	return endpoint.GetUserRequest{UID: value}, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
