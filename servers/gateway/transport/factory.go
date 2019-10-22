package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	kitTransport "github.com/go-kit/kit/transport/http"
)

func svcFactory(method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc kitTransport.EncodeRequestFunc
			dec kitTransport.DecodeResponseFunc
		)

		if method == http.MethodGet {
			enc, dec = encodeGetRequest, decodeResponse
		} else {
			enc, dec = encodeJsonRequest, decodeResponse
		}
		return kitTransport.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

func encodeGetRequest(ctx context.Context, r *http.Request, request interface{}) error {
	data, ok := request.(common.RequestUrlParams)
	if ok {
		r.URL.Path = strings.Replace(r.URL.Path, "{param}", data.Param, -1)
	}

	// return setTokenToHeader(ctx, r)
	return nil
}

func decodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var res common.Response

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

func encodeJsonRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)

	return nil
}

// 设置token到请求头
func setTokenToHeader(ctx context.Context, r *http.Request) error {
	if v := ctx.Value(common.AuthHeaderKey); v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		r.Header.Set(common.AuthHeaderKey, common.ServerAuthKey+" "+string(b))
	}
	return nil
}
