package main

import (
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	etcd2 "github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/gateway/router"
	"github.com/Zhan9Yunhua/blog-svr/gateway/server"
	"github.com/Zhan9Yunhua/blog-svr/shared/db"
	"github.com/Zhan9Yunhua/blog-svr/shared/logger"
	"github.com/Zhan9Yunhua/blog-svr/shared/middleware"
	"github.com/Zhan9Yunhua/blog-svr/shared/session"
	"github.com/Zhan9Yunhua/blog-svr/shared/zipkin"
	zipkinMiddlewareHttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func main() {
	conf := config.GetConfig()
	lg := logger.NewLogger(conf.LogPath)

	etcdClient := etcd2.NewEtcd()

	pool := db.NewRedis(conf.Redis)
	session := session.NewStorage(pool)

	zipkinTracer := zipkin.NewZipkin(lg, conf.ZipkinAddr, conf.ServerAddr, "gateway_server")

	r := router.NewRouter(lg, zipkinTracer)
	{
		r.Service("/svc/user", etcdClient)
		r.Post("/svc/user/login")
		r.Post("/svc/user/register")
		r.Get("/svc/user/code")
		r.Get("/svc/user/{param}")
		r.JetGet("/svc/user/{param}", middleware.CookieMiddleware(session))
	}

	handler := zipkinMiddlewareHttp.NewServerMiddleware(
		zipkinTracer,
		zipkinMiddlewareHttp.SpanName("gateway"),
		zipkinMiddlewareHttp.TagResponseSize(true),
		zipkinMiddlewareHttp.ServerTags(map[string]string{
			"component": "gateway_server",
		}),
	)(r.R)

	server.RunServer(conf.ServerAddr, handler)
}
