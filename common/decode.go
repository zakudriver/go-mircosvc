package common

import (
	"context"
	"encoding/json"
	"net/http"

	kitTransport "github.com/go-kit/kit/transport/http"
)

func DecodeEmpty(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func DecodeEmptyHttpRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeJsonRequest(reqPtr interface{}) kitTransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		if err := json.NewDecoder(r.Body).Decode(reqPtr); err != nil {
			return nil, err
		}
		return reqPtr, nil
	}
}
