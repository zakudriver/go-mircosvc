package transport

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	userPb "github.com/kum0/blog-svr/pb/user"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
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
			encodeGRPCLoginResponse,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "Login",
				logger)))...,
		),
		sendCode: kitGrpcTransport.NewServer(
			endpoints.SendCodeEP,
			decodeGRPCSendCodeRequest,
			encodeGRPCSendCodeResponse,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "SendCode",
				logger)))...,
		),
	}
}

type grpcServer struct {
	getUser  kitGrpcTransport.Handler `json:""`
	login    kitGrpcTransport.Handler `json:""`
	sendCode kitGrpcTransport.Handler `json:""`
}

func (gs *grpcServer) GetUser(ctx context.Context, req *userPb.GetUserRequest) (*userPb.GetUserReply, error) {
	_, rp, err := gs.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}
	rep, ok := rp.(*userPb.GetUserReply)
	if !ok {
		return nil, errors.New("*userPb.GetUserReply")
	}
	return &userPb.GetUserReply{Uid: rep.Uid}, nil
}

func (gs *grpcServer) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginReply, error) {
	_, rp, err := gs.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}

	rep, ok := rp.(*userPb.LoginReply)
	if !ok {
		return nil, errors.New("*userPb.LoginReply")
	}
	return rep, nil
}

func (gs *grpcServer) SendCode(ctx context.Context, req *userPb.SendCodeRequest) (*userPb.SendCodeReply, error) {
	_, rp, err := gs.sendCode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}

	rep, ok := rp.(*userPb.SendCodeReply)
	if !ok {
		return nil, errors.New("*userPb.SendCodeReply")
	}
	return rep, nil
}
