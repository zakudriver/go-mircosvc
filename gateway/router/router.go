package router

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"time"
)

func NewRouter(logger log.Logger) *Router {
	var router Router
	router.R= mux.NewRouter()
	router.Logger = logger
	return &router
}

type Router struct {
	R       *mux.Router
	EtcdIns *etcdv3.Instancer
	Logger  log.Logger
}

func (r *Router) Service(prefix string, etcdClient etcdv3.Client) {
	r.EtcdIns = etcd.GetEtcdIns(prefix, etcdClient, r.Logger)
}

func (r *Router) Post(path string, middlewares ...endpoint.Middleware) {
	r.R.Handle(path, MakeHandler(
		r.Logger,
		r.EtcdIns,
		"POST",
		path,
		middlewares...,
	))
}

func (r *Router) Get(path string, middlewares ...endpoint.Middleware) {
	r.R.Handle(path, MakeHandler(
		r.Logger,
		r.EtcdIns,
		"GET",
		path,
		middlewares...,
	))
}

func MakeHandler(logger log.Logger, ins *etcdv3.Instancer, method string, path string,
	middlewares ...endpoint.Middleware) *kithttp.Server {
	factory := SvcFactory(method, path)

	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	for _, middleware := range middlewares {
		retry = middleware(retry)
	}

	opts := []kithttp.ServerOption{
		// kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	decode := HandleDecodeRequest(method)
	return kithttp.NewServer(retry, decode, EncodeJSONResponse, opts...)
}
