package transport

import (
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
)

func MakeGRPCClient(conn *grpc.ClientConn, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) endpoints.IUserService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	options := []kitGrpcTransport.ClientOption{
		kitZipkin.GRPCClientTrace(zipkinTracer),
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.Usersvc",
			"GetUser",
			encodeGRPCGetUserRequest,
			decodeGRPCGetUserResponse,
			userPb.GetUserReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		getUserEndpoint = kitOpentracing.TraceClient(otTracer, "GetUser")(getUserEndpoint)
		getUserEndpoint = limiter(getUserEndpoint)
	}

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.Usersvc",
			"Login",
			encodeGRPCLoginRequest,
			decodeGRPCGetUserResponse,
			userPb.LoginReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		loginEndpoint = kitOpentracing.TraceClient(otTracer, "Login")(loginEndpoint)
		loginEndpoint = limiter(loginEndpoint)
	}

	return &endpoints.Endponits{
		GetUserEP: getUserEndpoint,
		LoginEP:   loginEndpoint,
	}
}

func grpcEncodeError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		return status.Error(st.Code(), st.Message())
	}

	return status.Error(codes.Internal, "internal server error")
}
