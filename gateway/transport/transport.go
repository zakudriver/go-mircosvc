package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(logger log.Logger, ins *etcdv3.Instancer, method string, path string,
	middlewares ...endpoint.Middleware) *kithttp.Server {
	factory := svcFactory(method, path)

	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	for _, m := range middlewares {
		retry = m(retry)
	}

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	var decode kithttp.DecodeRequestFunc
	if method == "GET" {
		decode = decodeGetRequest
	} else {
		decode = decodeJsonRequest
	}

	return kithttp.NewServer(retry, decode, encodeJsonResponse, opts...)
}

func encodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// 内部 -> 外部：解码get参数
func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		return nil, nil
	}

	value, err := vars["param"]
	if !err {
		return nil, common.ErrRouteArgs
	}

	var param common.RequestUrlParams
	param.Param = value

	return param, nil
}

// 内部 -> 外部 解析请求参数
func decodeJsonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request common.RequestBodyParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return request, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  err.Error(),
	})
}
