package transport

import (
	"context"
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
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/service"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
)

type grpcServer struct {
	getUser kitGrpcTransport.Handler `json:""`
}

func (s *grpcServer) GetUser(ctx context.Context, req *userPb.GetUserRequest) (*userPb.GetUserReply, error) {
	_, rp, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}
	rep := rp.(*userPb.GetUserReply)
	return &userPb.GetUserReply{Uid: rep.Uid}, nil
}

func NewGRPCClient(conn *grpc.ClientConn, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) service.IUserService {
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

	return endpoints.Endponits{
		GetUserEP: getUserEndpoint,
	}
}

func MakeGRPCServer(endpoints endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) userPb.UsersvcServer {

	options := []kitGrpcTransport.ServerOption{
		kitZipkin.GRPCServerTrace(zipkinTracer),
	}

	return &grpcServer{
		getUser: kitGrpcTransport.NewServer(
			endpoints.GetUserEP,
			decodeGRPCGetUserRequest,
			encodeGRPCGetUserResponse,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "GetUser",
				logger)))...,
		),
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
	switch err {
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
