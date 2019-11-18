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
	"github.com/kum0/blog-svr/shared/middleware"
	"github.com/kum0/blog-svr/shared/session"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"time"
)

func MakeHandler(
	etcdClient etcdv3.Client,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger,
	retryMax int,
	retryTimeout int,
	session *session.Storage,
) http.Handler {
	r := mux.NewRouter()

	// user endpoint
	{
		endpoints := new(usersvcEndpoints.Endponits)
		ins := sharedEtcd.NewInstancer("/usersvc", etcdClient, logger)
		{
			factory := usersvcFactory(usersvcEndpoints.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpoints.GetUserEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout, middleware.CookieMiddleware(session))
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeLoginEndpoint, tracer, zipkinTracer, logger)
			endpoints.LoginEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeRegisterEndpoint, tracer, zipkinTracer, logger)
			endpoints.RegisterEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeSendCodeEndpoint, tracer, zipkinTracer, logger)
			endpoints.SendCodeEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeUserListEndpoint, tracer, zipkinTracer, logger)
			endpoints.UserListEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeUserListEndpoint, tracer, zipkinTracer, logger)
			endpoints.UserListEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}

		{
			factory := usersvcFactory(usersvcEndpoints.MakeAuthEndpoint, tracer, zipkinTracer, logger)
			endpoints.AuthEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
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

func makeEndpoint(
	factory sd.Factory,
	ins *etcdv3.Instancer,
	logger log.Logger,
	retryMax int,
	retryTimeout int,
	middlewares ...endpoint.Middleware,
) endpoint.Endpoint {
	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	ep := lb.RetryWithCallback(time.Duration(retryTimeout)*time.Second, balancer, func(n int, received error) (bool,
		error) {
		if err := encodeError(received); err != nil {
			return false, err
		}
		return n < retryMax, nil
	})

	for _, m := range middlewares {
		ep = m(ep)
	}
	return ep
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
