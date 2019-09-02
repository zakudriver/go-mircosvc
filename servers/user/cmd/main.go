package main

import (
	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/server"
	"github.com/Zhan9Yunhua/blog-svr/services/db"
	"github.com/Zhan9Yunhua/blog-svr/services/email"
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

	conf := config.GetConfig()
	// mdb := db.NewMysql(conf.Mysql)
	rd := db.NewRedis(conf.Redis)
	email := email.NewEmail(conf.Email)

	var userSvc service.UserServicer
	userSvc = service.NewUserService(nil, rd, email)
	userSvc = middleware.InstrumentingMiddleware()(userSvc)

	mux := http.NewServeMux()
	mux.Handle(conf.Prefix+"/", service.MakeHandler(userSvc, lg))

	server.RunServer(mux, conf.ServerAddr)
}
