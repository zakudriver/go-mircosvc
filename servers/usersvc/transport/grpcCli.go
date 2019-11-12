package transport

import (
	"github.com/kum0/blog-svr/common"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	userPb "github.com/kum0/blog-svr/pb/user"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/time/rate"
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
		getUserEndpoint = limiter(getUserEndpoint)
		getUserEndpoint = kitOpentracing.TraceClient(otTracer, method)(getUserEndpoint)
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
			common.EncodeEmpty,
			decodeGRPCSendCodeResponse,
			userPb.SendCodeReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		sendCodeEndpoint = kitOpentracing.TraceClient(otTracer, method)(sendCodeEndpoint)
		sendCodeEndpoint = limiter(sendCodeEndpoint)
	}

	var registerEndpoint endpoint.Endpoint
	{
		method := "Register"
		registerEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.Usersvc",
			method,
			encodeGRPCRegisterRequest,
			decodeGRPCRegisterResponse,
			userPb.RegisterReply{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		registerEndpoint = kitOpentracing.TraceClient(otTracer, method)(registerEndpoint)
		registerEndpoint = limiter(registerEndpoint)
	}

	return &endpoints.Endponits{
		GetUserEP:  getUserEndpoint,
		LoginEP:    loginEndpoint,
		SendCodeEP: sendCodeEndpoint,
		RegisterEP: registerEndpoint,
	}
}

