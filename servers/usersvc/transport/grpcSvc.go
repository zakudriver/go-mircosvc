package transport

import (
	"context"
	"errors"
	"github.com/kum0/blog-svr/common"
	"google.golang.org/grpc/status"

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
			common.EncodeGRPCResponse(new(userPb.GetUserReply)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "GetUser",
				logger)))...,
		),
		login: kitGrpcTransport.NewServer(
			endpoints.LoginEP,
			decodeGRPCLoginRequest,
			common.EncodeGRPCResponse(new(userPb.LoginReply)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "Login",
				logger)))...,
		),
		sendCode: kitGrpcTransport.NewServer(
			endpoints.SendCodeEP,
			common.DecodeEmpty,
			common.EncodeGRPCResponse(new(userPb.SendCodeReply)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "SendCode",
				logger)))...,
		),
		register: kitGrpcTransport.NewServer(
			endpoints.RegisterEP,
			decodeGRPCRegisterRequest,
			common.EncodeEmpty,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "Register",
				logger)))...,
		),
	}
}

type grpcServer struct {
	getUser  kitGrpcTransport.Handler `json:""`
	login    kitGrpcTransport.Handler `json:""`
	sendCode kitGrpcTransport.Handler `json:""`
	register kitGrpcTransport.Handler `json:""`
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

func (gs *grpcServer) Register(ctx context.Context, req *userPb.RegisterRequest) (*userPb.RegisterReply, error) {
	_, _, err := gs.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(err)
	}

	return new(userPb.RegisterReply), nil
}

func grpcEncodeError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		// return status.Error(st.Code(), st.Message())
		return errors.New(st.Message())
	}

	// return status.Error(codes.Internal, "internal server error")
	return err
}
