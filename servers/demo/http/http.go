package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zhan9Yunhua/blog-svr/servers/demo/endpoint"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]kithttp.ServerOption) http.Handler {
	svr := kithttp.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse,
		options["Create"]...)

	return handleRouter(svr)
}

type responseWrapper struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// 空json返回
type nop struct{}

func handleRouter(svr *kithttp.Server) *mux.Router {
	m := mux.NewRouter()
	{
		handle := handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
			handlers.AllowedMethods([]string{"POST", "GET"}),
		)(svr)
		m.Methods("POST", "GET").Path("/order/create").Handler(handle)
	}

	return m
}

// 解析请求参数
func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.New("decodeCreateRequest")
	}

	return request, nil
}

// 编码响应
func encodeCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(kitendpoint.Failer); ok && nil != f.Failed() {
		return f.Failed()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// 成功的响应
	json.NewEncoder(w).Encode(responseWrapper{
		Code:    0,
		Message: "suc",
		Data:    response,
	})

	return nil
}

// 失败的统一格式化
func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	var code int
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err2code(err))

	json.NewEncoder(w).Encode(responseWrapper{
		Code:    code,
		Message: err.Error(),
		Data:    nop{},
	})
}

// 响应码设置
func err2code(err error) int {
	return http.StatusInternalServerError
}
