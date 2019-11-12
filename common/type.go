package common

// body参数
type RequestBodyParams = map[string]interface{}

// url参数
type RequestUrlParams struct {
	Param string `json:"param"`
}

// 响应格式
type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
