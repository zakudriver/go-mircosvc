package transport

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"github.com/kum0/blog-svr/common"
	userPb "github.com/kum0/blog-svr/pb/user"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

func MakeGRPCServer(eps *endpoints.Endponits, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) userPb.UsersvcServer {

	options := []kitGrpcTransport.ServerOption{
		kitZipkin.GRPCServerTrace(zipkinTracer),
	}

	return &grpcServer{
		getUser: kitGrpcTransport.NewServer(
			eps.GetUserEP,
			decodeGRPCGetUserRequest,
			common.EncodeGRPCResponse(new(userPb.GetUserResponse)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "GetUser",
				logger)))...,
		),
		login: kitGrpcTransport.NewServer(
			eps.LoginEP,
			decodeGRPCLoginRequest,
			common.EncodeGRPCResponse(new(userPb.LoginResponse)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "Login",
				logger)))...,
		),
		sendCode: kitGrpcTransport.NewServer(
			eps.SendCodeEP,
			common.DecodeEmpty,
			common.EncodeGRPCResponse(new(userPb.SendCodeResponse)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "SendCode",
				logger)))...,
		),
		register: kitGrpcTransport.NewServer(
			eps.RegisterEP,
			decodeGRPCRegisterRequest,
			common.EncodeEmpty,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "Register",
				logger)))...,
		),
		userList: kitGrpcTransport.NewServer(
			eps.UserListEP,
			decodeGRPCUserListRequest,
			encodeGRPCUserListResponse,
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "UserList",
				logger)))...,
		),
	}
}

type grpcServer struct {
	getUser  kitGrpcTransport.Handler `json:""`
	login    kitGrpcTransport.Handler `json:""`
	sendCode kitGrpcTransport.Handler `json:""`
	register kitGrpcTransport.Handler `json:""`
	userList kitGrpcTransport.Handler `json:""`
}

func (gs *grpcServer) GetUser(ctx context.Context, req *userPb.GetUserRequest) (*userPb.GetUserResponse, error) {
	_, rp, err := gs.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	rep, ok := rp.(*userPb.GetUserResponse)
	if !ok {
		return nil, errors.New("*userPb.GetUserResponse")
	}
	return &userPb.GetUserResponse{Uid: rep.Uid}, nil
}

func (gs *grpcServer) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	_, rp, err := gs.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	rep, ok := rp.(*userPb.LoginResponse)
	if !ok {
		return nil, errors.New("*userPb.LoginResponse")
	}
	return rep, nil
}

func (gs *grpcServer) SendCode(ctx context.Context, req *userPb.SendCodeRequest) (*userPb.SendCodeResponse, error) {
	_, rp, err := gs.sendCode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	rep, ok := rp.(*userPb.SendCodeResponse)
	if !ok {
		return nil, errors.New("*userPb.SendCodeResponse")
	}
	return rep, nil
}

func (gs *grpcServer) Register(ctx context.Context, req *userPb.RegisterRequest) (*userPb.RegisterResponse, error) {
	_, _, err := gs.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return new(userPb.RegisterResponse), nil
}

func (gs *grpcServer) UserList(ctx context.Context, req *userPb.UserListRequest) (*userPb.UserListResponse, error) {
	_, rp, err := gs.userList.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	rep, ok := rp.(*userPb.UserListResponse)
	if !ok {
		return nil, errors.New("*userPb.UserListResponse")
	}
	return rep, nil
}
