package common

import (
	"context"
	"encoding/json"
	"net/http"
)

func EncodeEmpty(_ context.Context, a interface{}) (request interface{}, err error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
