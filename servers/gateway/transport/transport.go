package transport

import (
	"io"
	"net/http"
	"time"

	usersvcEndpoints "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	usersvcTransport "github.com/Zhan9Yunhua/blog-svr/servers/usersvc/transport"
	"github.com/go-kit/kit/sd/lb"

	sharedEtcd "github.com/Zhan9Yunhua/blog-svr/shared/etcd"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
)

func MakeHandler(etcdClient etcdv3.Client, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// user endpoint
	{
		endpoints := &usersvcEndpoints.Endponits{}
		ins := sharedEtcd.NewInstancer("/usersvc", etcdClient, logger)
		{
			factory := usersvcFactory(usersvcEndpoints.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpointer := sd.NewEndpointer(ins, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(3, 3*time.Second, balancer)
			endpoints.GetUserEP = retry
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeLoginEndpoint, tracer, zipkinTracer, logger)
			endpointer := sd.NewEndpointer(ins, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(3, 3*time.Second, balancer)
			endpoints.LoginEP= retry
		}
		r.PathPrefix("/usersvc").Handler(http.StripPrefix("/usersvc", usersvcTransport.NewHTTPHandler(endpoints, tracer,
			zipkinTracer, logger)))
	}

	// article endpoint
	{
	}

	return r
}

func usersvcFactory(makeEndpoint func(service usersvcEndpoints.IUserService) endpoint.Endpoint, tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := usersvcTransport.MakeGRPCClient(conn, tracer, zipkinTracer, logger)
		return makeEndpoint(service), conn, nil
	}
}
