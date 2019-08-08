package main

import (
	discovery "github.com/Zhan9Yunhua/blog-svr/servers/lorem-etcd/etcd-sample"
)

func main() {
	worker := discovery.NewWorker("node-01", "127.0.0.1", discovery.Endpoints)
	worker.HeartBeat()
}
