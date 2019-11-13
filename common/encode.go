package common

import (
	"context"
	"encoding/json"
	kitTransportGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/kum0/blog-svr/utils"
	"net/http"
)

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(Response{
		Msg: err.Error(),
	})
}

func EncodeEmpty(_ context.Context, a interface{}) (request interface{}, err error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeGRPCResponse(a interface{}) kitTransportGrpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		// res, ok := response.(Response)
		// if !ok {
		// 	return nil, errors.New("encodeGRPCResponse: interface conversion error")
		// }

		if err := utils.StructCopy(response, a); err != nil {
			return nil, err
		}
		return a, nil
	}
}
