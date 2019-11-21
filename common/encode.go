package common

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitTransportGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/kum0/go-mircosvc/utils"
)

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	c, e := decodeError(err)

	w.WriteHeader(c)
	json.NewEncoder(w).Encode(Response{
		Msg: e.Error(),
	})
}

func EncodeEmpty(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// code := http.StatusOK
	// if sc, ok := response.(kitTransport.StatusCoder); ok {
	// 	code = sc.StatusCode()
	// }
	// w.WriteHeader(code)
	//
	// if code == http.StatusNoContent {
	// 	return nil
	// }

	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		return json.NewEncoder(w).Encode(Response{Msg: f.Failed().Error()})
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
