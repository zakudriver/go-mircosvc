package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	factory := SvcFactory(method, path)

	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	for _, middleware := range middlewares {
		retry = middleware(retry)
	}

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	decode := handleDecodeRequest(method)
	return kithttp.NewServer(retry, decode, EncodeJsonResponse, opts...)
}


func handleDecodeRequest(method string) kithttp.DecodeRequestFunc {
	if method == "GET" {
		return DecodeGetRequest
	}
	return DecodeJsonRequest

}


// 客户端到内部服务：转换Json响应
func EncodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}


// 内部服务到客户端：解码Get请求
func DecodeGetRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	fmt.Println("DecodeGetRequest")
	// vars := mux.Vars(req)
	// param, err := vars["param"]
	//
	// if !err {
	// 	return nil, errBadRoute
	// }
	// var getReq commonUrlReq
	// getReq.Param = param
	// return getReq, nil
	return nil, nil
}

// 内部服务到客户端：解码Json请求
func DecodeJsonRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	fmt.Println("DecodeJsonRequest")
	var request commonJsonReq
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return req, nil
}

type commonJsonReq struct {
	Param map[string]interface{} `json:"param"`
}

type commonUrlReq struct {
	Param string `json:"param"`
}

type commonRes struct {
	Code int                    `json:"code,string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
	Err  string                 `json:"err,omitempty"`
}
type outputRes struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// 错误码
var errBadRoute = errors.New("10010 错误的路由参数")

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  err.Error(),
	})
}
