package service

import (
	"context"
	"errors"

	"github.com/Zhan9Yunhua/blog-svr/servers/demo/entity"
)

type Middleware func(OrderService) OrderService

// 服务抽象
type OrderService interface {
	Create(ctx context.Context, orderId string) (entity.Order, error)
}

// 订单结构
type baseOrderService struct{}

// 创建订单
func (os baseOrderService) Create(ctx context.Context, orderId string) (o entity.Order, err error) {
	if "" == orderId {
		return o, errors.New("orderId")
	}

	o = entity.Order{

		Id:     "#" + orderId,
		Source: "APP",
		IsPay:  1,
	}

	return o, nil
}

// 服务对象实例化，并且组装中间件
func New(middleware []Middleware) OrderService {
	var svc = getBaseService()

	for _, m := range middleware {
		svc = m(svc)
	}

	return svc
}

// 获取当前实例
func getBaseService() OrderService {
	return &baseOrderService{}
}
