package main

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/gateway/logger"
	"github.com/Zhan9Yunhua/blog-svr/gateway/middleware"
	"github.com/Zhan9Yunhua/blog-svr/gateway/router"
	"github.com/Zhan9Yunhua/blog-svr/gateway/server"
	"github.com/Zhan9Yunhua/blog-svr/services/db"
	"github.com/Zhan9Yunhua/blog-svr/services/session"
)

func main() {
	conf := config.GetConfig()
	lg := logger.NewLogger()

	etcdClient := etcd.NewEtcd()

	pool := db.NewRedis(conf.Redis)
	session := session.NewStorage(pool)

	r := router.NewRouter(lg)
	{
		r.Service("/svc/user", etcdClient)
		r.Post("/svc/user/login")
		r.Post("/svc/user/register")
		r.JetGet("/svc/user/{param}", middleware.CookieMiddleware(session))
	}

	server.RunServer(config.GetConfig().ServerPort, r)
}
