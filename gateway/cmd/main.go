package main

import (
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/gateway/etcd"
	"github.com/Zhan9Yunhua/blog-svr/gateway/logger"
	"github.com/Zhan9Yunhua/blog-svr/gateway/router"
	"github.com/Zhan9Yunhua/blog-svr/gateway/server"
)

func main() {
	lg := logger.NewLogger()

	etcdClient := etcd.NewEtcd()

	fmt.Println(lg, etcdClient)

	r := router.NewRouter(lg)
	{
		r.Service("/svc/user", etcdClient)
		r.Get("/svc/user")
	}

	server.RunServer(config.GetConfig().ServerPort, r)
}
