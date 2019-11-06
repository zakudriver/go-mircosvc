package transport

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	userPb "github.com/kum0/blog-svr/pb/user"
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
		method := "GetUser"
		getUserEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.Usersvc",
			method,
			encodeGRPCGetUserRequest,
			decodeGRPCGetUserResponse,
			userPb.GetUserReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		getUserEndpoint = kitOpentracing.TraceClient(otTracer, method)(getUserEndpoint)
		getUserEndpoint = limiter(getUserEndpoint)
	}

	var loginEndpoint endpoint.Endpoint
	{
		method := "Login"
		loginEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.Usersvc",
			method,
			encodeGRPCLoginRequest,
			decodeGRPCLoginResponse,
			userPb.LoginReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		loginEndpoint = kitOpentracing.TraceClient(otTracer, method)(loginEndpoint)
		loginEndpoint = limiter(loginEndpoint)
	}

	var sendCodeEndpoint endpoint.Endpoint
	{
		method := "SendCode"
		sendCodeEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.Usersvc",
			method,
			encodeGRPCSendCodeRequest,
			decodeGRPCSendCodeResponse,
			userPb.SendCodeReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		sendCodeEndpoint = kitOpentracing.TraceClient(otTracer, method)(sendCodeEndpoint)
		sendCodeEndpoint = limiter(sendCodeEndpoint)
	}

	return &endpoints.Endponits{
		GetUserEP:  getUserEndpoint,
		LoginEP:    loginEndpoint,
		SendCodeEP: sendCodeEndpoint,
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
