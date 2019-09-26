package main

import (
	"net/http"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/shared/db"
	"github.com/Zhan9Yunhua/blog-svr/shared/email"
	"github.com/Zhan9Yunhua/blog-svr/shared/server"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/middleware"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/service"
	"github.com/Zhan9Yunhua/blog-svr/shared/etcd"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
)

func main() {
	conf := config.GetConfig()

	lg := logger.NewLogger(conf.LogPath)

	etcdClient := etcd.NewEtcd(conf.EtcdAddr)

	register := etcd.Register(conf.Prefix, conf.ServerAddr, etcdClient, lg)
	defer register.Deregister()

	var userSvc service.IUserService
	{
		mdb := db.NewMysql(conf.Mysql)
		rd := db.NewRedis(conf.Redis)
		email := email.NewEmail(conf.Email)

		userSvc = service.NewUserService(mdb, rd, email)
		userSvc = middleware.NewPrometheusMiddleware()(userSvc)
	}

	mux := http.NewServeMux()
	mux.Handle(conf.Prefix+"/", service.MakeHandler(userSvc, lg))

	server.RunServer(mux, conf.ServerAddr)
}
