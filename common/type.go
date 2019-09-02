package common

// body参数
type RequestBodyParams = map[string]interface{}

// url参数
type RequestUrlParams struct {
	Param string `json:"param"`
}

// 内部响应
type InnerResponse struct {
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
	Err  string                 `json:"err,omitempty"`
}

// 外部输出响应
type Response struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
