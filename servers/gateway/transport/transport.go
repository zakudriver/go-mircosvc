package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	usersvcEndpoints "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	usersvcSer "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	usersvcTransport "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"time"
)

func MakeHandler(ctx context.Context, etcdClient etcdv3.Client, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	// conf := config.GetConfig()
	r := mux.NewRouter()
	{
		endpoints := usersvcEndpoints.Endponits{}
		// ins, err := etcdv3.NewInstancer(etcdClient, "/usersvc", logger)
		// if err != nil {
		// 	logger.Log(err)
		// }

		{
			factory, _ := usersvcfactory("localhost:5002", usersvcEndpoints.MakeGetUserEndpoint, tracer,
				zipkinTracer,
				logger)
			// endpointer := sd.NewEndpointer(ins, factory, logger)
			// balancer := lb.NewRoundRobin(endpointer)

			// retry := lb.Retry(conf.RetryMax, time.Duration(conf.RetryTimeout), balancer)
			endpoints.GetUserEP = factory
		}
		r.PathPrefix("/usersvc").Handler(http.StripPrefix("/usersvc", usersvcTransport.NewHTTPHandler(endpoints, tracer,
			zipkinTracer, logger)))
	}

	return r
}

func usersvcfactory(
	addr string,
	makeEndpoint func(service usersvcSer.IUserService) endpoint.Endpoint,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) (endpoint.Endpoint, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	svc := usersvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)

	return makeEndpoint(svc), nil
}

func usersvcFactory(makeEndpoint func(service usersvcSer.IUserService) endpoint.Endpoint, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := usersvcTransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}

func MakeHandler(logger log.Logger, ins *etcdv3.Instancer, method string, path string, isCookie bool,
	middlewares ...endpoint.Middleware) *kithttp.Server {
	factory := svcFactory(method, path)

	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	for _, m := range middlewares {
		retry = m(retry)
	}

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		// kithttp.ServerFinalizer(func(ctx context.Context, code int, r *http.Request) {
		// }),
	}

	if isCookie {
		opts = append(opts, kithttp.ServerBefore(cookieToContext()))
	}

	var decode kithttp.DecodeRequestFunc
	if method == "GET" {
		decode = decodeGetRequest
	} else {
		decode = decodeJsonRequest
	}

	return kithttp.NewServer(retry, decode, encodeJsonResponse, opts...)
}

func encodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// 内部 -> 外部：解码get参数
func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		return nil, nil
	}

	value, err := vars["param"]
	if !err {
		return nil, common.ErrRouteArgs
	}

	var param common.RequestUrlParams
	param.Param = value

	return param, nil
}

// 内部 -> 外部 解析请求参数
func decodeJsonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request common.RequestBodyParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return request, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(common.Response{
		Code: common.Error.Code(),
		Msg:  err.Error(),
	})
}

func cookieToContext() kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		c, err := r.Cookie(common.AuthHeaderKey)
		if err != nil {
			return ctx
		}

		return context.WithValue(ctx, common.SessionKey, c.Value)
	}
}
