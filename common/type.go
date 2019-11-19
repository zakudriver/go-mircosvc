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
	Err  error       `json:"-"`
	Code int         `json:"-"`
}

func (r *Response) Failed() error {
	return r.Err
}

func (r *Response) StatusCode() int {
	if r.Code == 0 {
		return 200
	}
	return r.Code
}
