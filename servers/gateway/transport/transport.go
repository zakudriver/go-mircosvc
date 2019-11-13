package transport

import (
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/gorilla/mux"
	usersvcEndpoints "github.com/kum0/blog-svr/servers/usersvc/endpoints"
	usersvcTransport "github.com/kum0/blog-svr/servers/usersvc/transport"
	sharedEtcd "github.com/kum0/blog-svr/shared/etcd"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"time"
)

func MakeHandler(etcdClient etcdv3.Client, tracer opentracing.Tracer, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// user endpoint
	{
		endpoints := new(usersvcEndpoints.Endponits)
		ins := sharedEtcd.NewInstancer("/usersvc", etcdClient, logger)
		{
			factory := usersvcFactory(usersvcEndpoints.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpoints.GetUserEP = makeEndpoint(factory, ins, logger)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeLoginEndpoint, tracer, zipkinTracer, logger)
			endpoints.LoginEP = makeEndpoint(factory, ins, logger)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeRegisterEndpoint, tracer, zipkinTracer, logger)
			endpoints.RegisterEP = makeEndpoint(factory, ins, logger)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeSendCodeEndpoint, tracer, zipkinTracer, logger)
			endpoints.SendCodeEP = makeEndpoint(factory, ins, logger)
		}
		r.PathPrefix("/usersvc").Handler(http.StripPrefix("/usersvc", usersvcTransport.MakeHTTPHandler(endpoints, tracer,
			zipkinTracer, logger)))
	}

	// article endpoint
	{
	}

	return r
}

func usersvcFactory(
	makeEndpoint func(service usersvcEndpoints.IUserService) endpoint.Endpoint,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger,
) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := usersvcTransport.MakeGRPCClient(conn, tracer, zipkinTracer, logger)
		return makeEndpoint(service), conn, nil
	}
}

func makeEndpoint(factory sd.Factory, ins *etcdv3.Instancer, logger log.Logger) endpoint.Endpoint {
	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	return lb.RetryWithCallback(3*time.Second, balancer, func(n int, received error) (keepTrying bool, replacement error) {
		if err := encodeError(received); err != nil {
			return false, err
		}
		return n < 3, nil
	})
}

func encodeError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		if st.Code() == codes.InvalidArgument {
			return errors.New(st.Message())
		}
	}

	return nil
}
