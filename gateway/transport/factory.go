package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Zhan9Yunhua/logger"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
)

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
	data, ok := request.(commonUrlReq)
	fmt.Printf("%+v\n", data)
	if ok {
		req.URL.Path = strings.Replace(req.URL.Path, "{param}", data.Param, -1)
	}

	return nil
}

// 客户端到内部服务：转换Json请求
func EncodeJsonRequest(_ context.Context, req *http.Request, request interface{}) error {
	fmt.Println("EncodeJsonRequest")
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)

	return nil
}

// 内部服务到客户端：解码Get响应
func DecodeGetResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	fmt.Println("DecodeGetResponse")
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
	return resp, nil
}
