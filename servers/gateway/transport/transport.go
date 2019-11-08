package transport

import (
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/sd/lb"
	usersvcEndpoints "github.com/kum0/blog-svr/servers/usersvc/endpoints"
	usersvcTransport "github.com/kum0/blog-svr/servers/usersvc/transport"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/gorilla/mux"
	sharedEtcd "github.com/kum0/blog-svr/shared/etcd"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	zipkinGrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
)

func MakeHandler(etcdClient etcdv3.Client, tracer opentracing.Tracer, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// user endpoint
	{
		endpoints := &usersvcEndpoints.Endponits{}
		ins := sharedEtcd.NewInstancer("/usersvc", etcdClient, logger)
		{
			factory := usersvcFactory(usersvcEndpoints.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpoints.GetUserEP = makeEndpoint(factory, zipkinTracer, "GetUser", ins, logger)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeLoginEndpoint, tracer, zipkinTracer, logger)
			endpoints.LoginEP = makeEndpoint(factory, zipkinTracer, "Login", ins, logger)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeSendCodeEndpoint, tracer, zipkinTracer, logger)
			endpoints.SendCodeEP = makeEndpoint(factory, zipkinTracer, "SenCode", ins, logger)
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
		conn, err := grpc.Dial(instance, grpc.WithInsecure(), grpc.WithStatsHandler(zipkinGrpc.NewClientHandler(
			zipkinTracer)))
		if err != nil {
			return nil, nil, err
		}
		service := usersvcTransport.MakeGRPCClient(conn, tracer, zipkinTracer, logger)
		return makeEndpoint(service), conn, nil
	}
}

func makeEndpoint(factory sd.Factory, zipkinTracer *zipkin.Tracer, method string, ins *etcdv3.Instancer,
	logger log.Logger) endpoint.Endpoint {
	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	ep := lb.Retry(3, 3*time.Second, balancer)
	// return kitZipkin.TraceEndpoint(zipkinTracer, method)(ep)
	return ep
}
