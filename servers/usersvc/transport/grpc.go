package transport

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/addsvc/pb"
	"github.com/go-kit/kit/ratelimit"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	kitGrpcTransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"time"

	"github.com/opentracing/opentracing-go"
)

func NewGRPCClient(conn *grpc.ClientConn, otTracer opentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.AddsvcService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	zipkinClient := zipkin.GRPCClientTrace(zipkinTracer)

	options := []kitGrpcTransport.ClientOption{
		zipkinClient,
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = kitGrpcTransport.NewClient(
			conn,
			"pb.UserSvc",
			"GetUser",
			encodeGRPCSumRequest,
			decodeGRPCSumResponse,
			pb.SumReply{},
			append(options, grpctransport.ClientBefore(kitOpentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		sumEndpoint = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
		sumEndpoint = limiter(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Sum",
			Timeout: 30 * time.Second,
		}))(sumEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = grpctransport.NewClient(
			conn,
			"pb.Addsvc",
			"Concat",
			encodeGRPCConcatRequest,
			decodeGRPCConcatResponse,
			pb.ConcatReply{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		concatEndpoint = opentracing.TraceClient(otTracer, "Concat")(concatEndpoint)
		concatEndpoint = limiter(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Concat",
			Timeout: 30 * time.Second,
		}))(concatEndpoint)
	}

	return endpoints.Endpoints{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}
