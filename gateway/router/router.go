package router

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/gateway/transport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
)

func NewRouter(logger log.Logger, tracer *zipkin.Tracer) *Router {
	return &Router{
		R:         mux.NewRouter(),
		EtcdIns:   nil,
		Logger:    logger,
		// Transport: zk.NewTransport(tracer),
	}
}

type Router struct {
	R         *mux.Router
	EtcdIns   *etcdv3.Instancer
	Logger    log.Logger
	// Transport http.RoundTripper
}

func (r *Router) Service(prefix string, etcdClient etcdv3.Client) {
	r.EtcdIns = etcd.GetEtcdIns(prefix, etcdClient, r.Logger)
}

func (r *Router) Post(path string, middlewares ...endpoint.Middleware) {
	r.R.Handle(path, transport.MakeHandler(
		r.Logger,
		r.EtcdIns,
		"POST",
		path,
		false,
		middlewares...,
	))
}

func (r *Router) JwtPost(path string, middlewares ...endpoint.Middleware) {
	r.R.Handle(path, transport.MakeHandler(
		r.Logger,
		r.EtcdIns,
		"POST",
		path,
		true,
		middlewares...,
	))
}

func (r *Router) Get(path string, middlewares ...endpoint.Middleware) {
	r.R.Handle(path, transport.MakeHandler(
		r.Logger,
		r.EtcdIns,
		"GET",
		path,
		false,
		middlewares...,
	))
}

func (r *Router) JetGet(path string, middlewares ...endpoint.Middleware) {
	r.R.Handle(path, transport.MakeHandler(
		r.Logger,
		r.EtcdIns,
		"GET",
		path,
		true,
		middlewares...,
	))
}
