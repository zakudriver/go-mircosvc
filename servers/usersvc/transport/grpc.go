package transport

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

type grpcServer struct {
	getUser kitGrpcTransport.Handler `json:""`
}

func (s *grpcServer) GetUser(ctx context.Context, req *userPb.GetUserRequest) (rep *userPb.GetUserReply, err error) {
	_, rp, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}
	rep = rp.(*userPb.GetUserReply)
	return rep, nil
}

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

func MakeGRPCServer(endpoints endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) (req userPb.UserSvcServer) {
	zipkinServer := kitZipkin.GRPCServerTrace(zipkinTracer)

	options := []kitGrpcTransport.ServerOption{
		kitGrpcTransport.ServerErrorLogger(logger),
		zipkinServer,
	}

	return &grpcServer{
		getUser: kitGrpcTransport.NewServer(
			endpoints.GetUserEP,
			decodeGRPCSumRequest,
			encodeGRPCSumResponse,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "GetUser",
				logger)))...,
		),
	}
}

func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	// req := grpcReq.(*pb.SumRequest)
	return endpoints.GetUserRequest{Uid: grpcReq.(string)}, nil
}

func encodeGRPCSumResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	// reply := grpcReply.(endpoints.SumResponse)
	return &userPb.GetUserReply{Uid: grpcReply.(string)}, nil;
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
