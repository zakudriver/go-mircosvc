package transport

import (
	"context"
	"errors"
	userPb "github.com/Zhan9Yunhua/blog-svr/pb/user"
	"github.com/Zhan9Yunhua/blog-svr/servers/usersvc/endpoints"
	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

func MakeGRPCServer(endpoints *endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
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
		login: kitGrpcTransport.NewServer(
			endpoints.LoginEP,
			decodeGRPCLoginRequest,
			encodeGRPCGetUserResponse,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "Login",
				logger)))...,
		),
	}
}

type grpcServer struct {
	getUser kitGrpcTransport.Handler `json:""`
	login   kitGrpcTransport.Handler `json:""`
}

func (s *grpcServer) GetUser(ctx context.Context, req *userPb.GetUserRequest) (*userPb.GetUserReply, error) {
	_, rp, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}
	rep := rp.(*userPb.GetUserReply)
	return &userPb.GetUserReply{Uid: rep.Uid}, nil
}

func (s *grpcServer) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginReply, error) {
	_, rp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}
	rep, ok := rp.(*userPb.LoginReply)
	if !ok {
		return nil, errors.New("*userPb.LoginReply")
	}
	return rep, nil
}
