package main

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	"github.com/Zhan9Yunhua/blog-svr/gateway/middleware"
	"github.com/Zhan9Yunhua/blog-svr/gateway/router"
	"github.com/Zhan9Yunhua/blog-svr/gateway/server"
	"github.com/Zhan9Yunhua/blog-svr/shared/db"
	"github.com/Zhan9Yunhua/blog-svr/shared/session"
)

func main() {
	conf := config.GetConfig()
	lg := logger.NewLogger(conf.LogPath)

	etcdClient := etcd.NewEtcd()

	pool := db.NewRedis(conf.Redis)
	session := session.NewStorage(pool)

	r := router.NewRouter(lg)
	{
		r.Service("/svc/user", etcdClient)
		r.Post("/svc/user/login")
		r.Post("/svc/user/register")
		r.Get("/svc/user/code")
		r.Get("/svc/user/{param}")
		r.JetGet("/svc/user/{param}", middleware.CookieMiddleware(session))
	}

	server.RunServer(conf.ServerPort, r)
}
