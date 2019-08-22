package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	value, ok := vars["Username"]
	if !ok {
		return nil, errBadRoute
	}
	fmt.Println("request: ", value)
	return loginRequest{Username: value}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// MakeHandler returns a handler for the booking service.
func MakeHandler(bs UcenterServiceInterface, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getUserHandler := kithttp.NewServer(
		makeGetUserEndpoint(bs),
		decodeLoginRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	conf := config.GetConfig()
	// 接口路由
	r.Handle(conf.Prefix+"/login", getUserHandler).Methods("POST")

	return r
}

// encode errors from business-logic
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
		"msg":  "from ucenter error: " + err.Error(),
	})
}

// var ErrUnknown = errors.New("unknown cargo")
// var ErrInvalidArgument = errors.New("error argument")
var errBadRoute = errors.New("error argument")