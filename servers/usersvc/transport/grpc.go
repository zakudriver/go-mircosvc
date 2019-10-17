package transport

import (
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"

	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"time"

	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	kitGrpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
)

func NewGRPCClient(conn *grpc.ClientConn, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) service.IUserService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	zipkinClient := kitZipkin.GRPCClientTrace(zipkinTracer)

	options := []kitGrpcTransport.ClientOption{
		zipkinClient,
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.UserSvc",
			"GetUser",
			decodeGRPCGetUserRequest,
			encodeGRPCGetUserResponse,
			userPb.GetUserReply{},
			append(options, kitGrpctransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		getUserEndpoint = kitOpentracing.TraceClient(otTracer, "GetUser")(getUserEndpoint)
		getUserEndpoint = limiter(getUserEndpoint)
		// getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		// 	Name:    "GetUser",
		// 	Timeout: 30 * time.Second,
		// }))(getUserEndpoint)
	}

	return endpoints.Endponits{
		GetUserEP: getUserEndpoint,
	}
}
