package transport

import (
	"context"
	"errors"

	articlePb "github.com/kum0/go-mircosvc/pb/article"
)

func decodeGRPCGetCategoriesResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	rp, ok := grpcResponse.(*articlePb.GetCategoriesResponse)
	if !ok {
		return nil, errors.New("decodeGRPCGetCategoriesResponse: interface conversion error")
	}

	return rp, nil
}
