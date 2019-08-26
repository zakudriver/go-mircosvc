package main

import (
	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/server"
	"net/http"

	_ "github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/etcd"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/logger"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/middleware"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/service"
)

func main() {
	lg := logger.NewLogger()

	etcdClient := etcd.NewEtcd()

	register := etcd.Register(etcdClient, lg)
	defer register.Deregister()

	var ucenterSvc service.UserServicer
	ucenterSvc = service.UserService{}
	ucenterSvc = middleware.InstrumentingMiddleware()(ucenterSvc)

	conf := config.GetConfig()

	mux := http.NewServeMux()
	mux.Handle(conf.Prefix+"/", service.MakeHandler(ucenterSvc, lg))

	server.RunServer(mux, conf.ServerAddr)
}
