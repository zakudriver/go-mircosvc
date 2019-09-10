package common

// body参数
type RequestBodyParams = map[string]interface{}

// url参数
type RequestUrlParams struct {
	Param string `json:"param"`
}

// 响应数据
type ResponseData = map[string]interface{}

// 相应格式
type Response struct {
	Code int32        `json:"code"`
	Msg  string       `json:"msg"`
	Data ResponseData `json:"data"`
}
