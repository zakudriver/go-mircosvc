package endpoint

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	"github.com/go-kit/kit/endpoint"
)

type Endponits struct {
	GetUserEP  endpoint.Endpoint
	LoginEP    endpoint.Endpoint
	SendCodeEP endpoint.Endpoint
}

func MakeGetUserEndpoint(svc service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)

		name, err := svc.GetUser(req.UID)
		if err != nil {
			return nil, err
		}
		data := map[string]interface{}{
			"id": name,
		}

		return common.Response{Code: common.OK.Code(), Msg: "ok", Data: data,}, nil
	}
}
