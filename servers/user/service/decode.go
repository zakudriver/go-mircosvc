package service

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/gorilla/mux"
	"net/http"
)

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	value, ok := vars["UID"]
	if !ok {
		return nil, common.ErrRouteArgs
	}
	return getUserRequest{UID: value}, nil
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
