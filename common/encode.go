package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitTransportGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/kum0/blog-svr/utils"
	"net/http"

	kitTransport "github.com/go-kit/kit/transport/http"
)

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Msg: err.Error(),
	})
}

func EncodeEmpty(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println(response)
	if f, ok := response.(endpoint.Failer); ok {
		fmt.Println(f.Failed().Error())
	}
	code := http.StatusOK
	if sc, ok := response.(kitTransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

func EncodeGRPCResponse(a interface{}) kitTransportGrpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		res, ok := response.(Response)
		if !ok {
			return nil, errors.New("encodeGRPCResponse: interface conversion error")
		}

		if err := utils.StructCopy(res.Data, a); err != nil {
			return nil, err
		}
		return a, nil
	}
}
