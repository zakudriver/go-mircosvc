package service

import (
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHandler(bs UserServicer, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	loginHandler := kithttp.NewServer(
		makeLoginEndpoint(bs),
		decodeLoginRequest,
		encodeResponseSetCookie,
		opts...,
	)

	getUserHandler := kithttp.NewServer(
		makeGetUserEndpoint(bs),
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	conf := config.GetConfig()
	// 接口路由
	r.Handle(conf.Prefix+"/login", loginHandler).Methods("POST")
	r.Handle(conf.Prefix+"/{UID}", getUserHandler).Methods("GET")

	return r
}

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

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// switch err {
	// 	case ErrUnknown:
	// 		w.WriteHeader(http.StatusNotFound)
	// 	case ErrInvalidArgument:
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	default:
	// 		w.WriteHeader(http.StatusInternalServerError)
	// }
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusNotFound,
		"msg":  "from user error: " + err.Error(),
	})
}
