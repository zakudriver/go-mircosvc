package transport

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"github.com/kum0/go-mircosvc/common"
	articlePb "github.com/kum0/go-mircosvc/pb/article"
	"github.com/kum0/go-mircosvc/servers/article/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

func MakeGRPCServer(eps *endpoints.Endpoints, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) articlePb.ArticlesvcServer {

	options := []kitGrpcTransport.ServerOption{
		kitZipkin.GRPCServerTrace(zipkinTracer),
	}

	return &grpcServer{
		getCategories: kitGrpcTransport.NewServer(
			eps.GetCategoriesEP,
			common.DecodeEmpty,
			common.EncodeGRPCResponse(new(articlePb.GetCategoriesResponse)),
			append(options, kitGrpcTransport.ServerBefore(kitOpentracing.GRPCToContext(otTracer, "GetCategories",
				logger)))...,
		),
	}
}

type grpcServer struct {
	getCategories kitGrpcTransport.Handler `json:""`
}

func (gs *grpcServer) GetCategories(ctx context.Context, req *articlePb.GetCategoriesRequest) (*articlePb.GetCategoriesResponse, error) {
	_, rp, err := gs.getCategories.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	rep, ok := rp.(*articlePb.GetCategoriesResponse)
	if !ok {
		return nil, errors.New("*articlePb.GetCategoriesResponse")
	}
	return rep, nil
}
