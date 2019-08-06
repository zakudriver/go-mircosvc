package main


import (
	"github.com/Zhan9Yunhua/blog-svr/servers/demo/endpoint"
	"github.com/Zhan9Yunhua/blog-svr/servers/demo/http"
	"github.com/Zhan9Yunhua/blog-svr/servers/demo/service"
	"log"
	nethttp "net/http"
	kithttp "github.com/go-kit/kit/transport/http"
)

func main() {
	httpHandler := createService()
	log.Println("demo1服务启动，服务地址：127.0.0.1:8088")
	err := nethttp.ListenAndServe(":8088", httpHandler)

	if nil != err {
		log.Println(err)
	}
}

// 创建服务
func createService() nethttp.Handler {
	// 创建业务对象
	svc := service.New(nil)
	// 创建端点对象
	eps := endpoint.New(svc, nil)
	// 设置http服务服务中间件
	options := defaultHttpOptions()
	// 端点绑定到http服务上
	return http.NewHTTPHandler(eps, options)
}

// HTTP服务中间件（服务的aop）
func defaultHttpOptions() map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{
		"Create": {
			kithttp.ServerErrorEncoder(http.ErrorEncoder),
		},
	}
	return options
}
