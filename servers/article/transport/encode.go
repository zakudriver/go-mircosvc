package transport

import (
	"context"
	"errors"

	"github.com/kum0/go-mircosvc/common"
	articlePb "github.com/kum0/go-mircosvc/pb/article"
)

func encodeGRPCGetCategoriesResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCGetCategoriesResponse: interface conversion error")
	}

	data := res.Data.(*articlePb.GetCategoriesResponse)

	return data, nil
}
