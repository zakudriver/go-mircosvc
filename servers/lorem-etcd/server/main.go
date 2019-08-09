package main

import (
	"context"
	"fmt"
	lorem_etcd "github.com/Zhan9Yunhua/blog-svr/servers/lorem-etcd"
	"net/http"
)

func main() {
	// var (
	// 	// 由于consul服务运行在docker或compose, 所以这两个地址一定要正确.
	// 	etcdURL       = os.Getenv("ETCD_URL")
	// 	etcdPrefix    = os.Getenv("ETCD_PREFIX")
	// 	advertiseAddr = os.Getenv("SERVER_ADDR")
	// 	advertisePort = os.Getenv("SERVER_PORT")
	// )

	var (
		etcdURL       = "127.0.0.1:2379"
		etcdPrefix    = "/test/"
		advertiseAddr = "127.0.0.1"
		advertisePort = "7999"
	)

	var svc lorem_etcd.Service
	svc = lorem_etcd.LoremService{}

	loremEndpoint := lorem_etcd.MakeLoremEndpoint(svc)
	endpoints := lorem_etcd.Endpoints{
		LoremEndpoint:  loremEndpoint,
	}

	client, err := lorem_etcd.ConnectEtcd(etcdURL)
	if err != nil {
		panic(err)
	}
	// 注册服务
	serviceAddr := advertiseAddr + ":" + advertisePort
	key := etcdPrefix + serviceAddr
	registrar := lorem_etcd.Register(client, key, serviceAddr)
	// 将go-kit类型的endpoint接口转换成http标准库接口
	registrar.Register()
	fmt.Println("register success")

	ctx := context.Background()
	handler := lorem_etcd.MakeHTTPHandler(ctx, endpoints)

	// 提供标准http服务
	fmt.Println("Starting server")
	fmt.Println(http.ListenAndServe(advertiseAddr+":"+advertisePort, handler))

	registrar.Deregister()
}
