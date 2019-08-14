package main

import (
	"github.com/Zhan9Yunhua/blog-svr/servers/user/etcd"
	_ "github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/logger"
)

func main() {
	lg := logger.NewLogger()


	etcdClient := etcd.NewEtcd()

	register:=etcd.Register(etcdClient,lg)
	defer register.Deregister()


}