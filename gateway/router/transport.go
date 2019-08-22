package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zhan9Yunhua/logger"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func HandleDecodeRequest(method string) kithttp.DecodeRequestFunc {
	if method == "GET" {
		return DecodeGetRequest
	}
	return DecodeJsonRequest

}

func SvcFactory(method string, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		tgt, err := url.Parse(instance)
		logger.Infoln("listening svc: ", tgt)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc kithttp.EncodeRequestFunc
			dec kithttp.DecodeResponseFunc
		)

		method = strings.ToUpper(method)

		if method == "GET" {
			enc, dec = EncodeGetRequest, DecodeGetResponse
		} else {
			enc, dec = EncodeJsonRequest, DecodeGetResponse
		}

		return kithttp.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

// 客户端到内部服务：转换Get请求
func EncodeGetRequest(_ context.Context, req *http.Request, request interface{}) error {
	fmt.Printf("%+v\n", request)
	data, ok := request.(commonUrlReq)
	fmt.Println(ok)
	if ok {
		req.URL.Path = strings.Replace(req.URL.Path, "{param}", data.Param, -1)
	}

	return nil
}

// 客户端到内部服务：转换Json请求
func EncodeJsonRequest(_ context.Context, req *http.Request, request interface{}) error {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)

	return nil
}

// 客户端到内部服务：转换Json响应
func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// 内部服务到客户端：解码Get响应
func DecodeGetResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	// var commonResponse commonRes
	// var outputResponse outputRes
	//
	// if err := json.NewDecoder(resp.Body).Decode(&commonResponse); err != nil {
	// 	return nil, err
	// }
	// if commonResponse.Err != "" {
	// 	outputResponse.Msg = commonResponse.Err
	// 	outputResponse.Code = 500
	// 	outputResponse.Data = map[string]interface{}{}
	// } else {
	// 	outputResponse.Msg = commonResponse.Msg
	// 	outputResponse.Code = commonResponse.Code
	// 	outputResponse.Data = commonResponse.Data
	// }
	// return outputResponse, nil
	return nil, nil
}

// 内部服务到客户端：解码Get请求
func DecodeGetRequest(ctx context.Context, req *http.Request) (interface{}, error) {
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
	var request commonJsonReq
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
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
