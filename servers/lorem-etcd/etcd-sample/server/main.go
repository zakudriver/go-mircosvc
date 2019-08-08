package main

import (
	discovery "github.com/Zhan9Yunhua/blog-svr/servers/lorem-etcd/etcd-sample"
)

func main() {
	master := discovery.NewMaster(discovery.Endpoints)
	master.WatchWorkers()
}
