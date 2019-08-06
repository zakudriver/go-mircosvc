package endpoint

import (
	"context"

	"github.com/Zhan9Yunhua/blog-svr/servers/demo/entity"
	"github.com/Zhan9Yunhua/blog-svr/servers/demo/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateEndpoint endpoint.Endpoint
}

// 封装打包多个端点，并使用包装模式加载中间件
func New(svc service.OrderService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateEndpoint: makeCreateEndpoint(svc),
	}

	for _, m := range mdw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}

	return eps
}


type CreateRequest struct {
	OrderId string `json:"orderId"`
}

// 响应参数格式
type CreateResponse struct {
	entity.Order
	err error
}

// create操作端点
func makeCreateEndpoint(svc service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		order, err := svc.Create(ctx, req.OrderId)

		return CreateResponse{
			order,
			err,
		}, nil
	}
}

// 错误获取
func (rs CreateResponse) Failed() error {
	return rs.err
}
