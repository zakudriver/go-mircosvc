package transport

import (
	"time"

	"github.com/kum0/go-mircosvc/common"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	articlePb "github.com/kum0/go-mircosvc/pb/article"
	"github.com/kum0/go-mircosvc/servers/article/endpoints"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func MakeGRPCClient(conn *grpc.ClientConn, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer,
	logger log.Logger) endpoints.ArticleServicer {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	options := []kitGrpcTransport.ClientOption{
		kitZipkin.GRPCClientTrace(zipkinTracer),
	}

	var GetCategories endpoint.Endpoint
	{
		method := "GetCategories"
		GetCategories = kitGrpcTransport.NewClient(
			conn,
			"pb.Articlesvc",
			method,
			common.EncodeEmpty,
			decodeGRPCGetCategoriesResponse,
			articlePb.GetCategoriesResponse{},
			append(options, kitGrpcTransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		GetCategories = limiter(GetCategories)
		GetCategories = kitOpentracing.TraceClient(otTracer, method)(GetCategories)
	}

	return &endpoints.Endpoints{
		GetCategoriesEP: GetCategories,
	}
}
