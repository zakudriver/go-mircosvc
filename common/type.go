package common

import "net/http"

// body参数
type RequestBodyParams = map[string]interface{}

// url参数
type RequestUrlParams struct {
	Param string `json:"param"`
}

// 响应格式
type Response struct {
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Err    error       `json:"-"`
	Header http.Header `json:"-"`
}

func (r Response) Failed() error {
	return r.Err
}

func (r Response) Headers() http.Header {
	return r.Header
}
