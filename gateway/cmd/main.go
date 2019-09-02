package main

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/gateway/logger"
	"github.com/Zhan9Yunhua/blog-svr/gateway/router"
	"github.com/Zhan9Yunhua/blog-svr/gateway/server"
)

func main() {
	conf := config.GetConfig()
	lg := logger.NewLogger()

	etcdClient := etcd.NewEtcd()

	// pool := db.NewRedis(conf.Redis)
	// session := session.NewStorage(pool)

	r := router.NewRouter(lg)
	{
		r.Service("/svc/user", etcdClient)
		r.Post("/svc/user/login")
		r.Post("/svc/user/register")
		r.Get("/svc/user/code")
		r.Get("/svc/user/{param}")
		// r.JetGet("/svc/user/{param}", middleware.CookieMiddleware(session))
	}

	server.RunServer(conf.ServerPort, r)
}
