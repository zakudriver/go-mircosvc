package router

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/gateway/transport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/gorilla/mux"
)

func NewRouter(logger log.Logger) *Router {
	var router Router
	router.R = mux.NewRouter()
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
