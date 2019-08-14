package etcd

import (
	"context"
	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/logger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"time"
)

func NewEtcd() etcdv3.Client {
	etcdClient, err := newEtcd()
	if err != nil {
		logger.Fatalln(err)
	}

	return etcdClient
}

func newEtcd() (etcdv3.Client, error) {
	etcdAddr := config.GetConfig().EtcdAddr

	options := etcdv3.ClientOptions{
		// Path to trusted ca file
		CACert: "",
		// Path to certificate
		Cert: "",
		// Path to private key
		Key: "",
		// Username if required
		Username: "",
		// Password if required
		Password: "",
		// If DialTimeout is 0, it defaults to 3s
		DialTimeout: time.Second * 3,
		// If DialKeepAlive is 0, it defaults to 3s
		DialKeepAlive: time.Second * 3,
	}

	ctx := context.Background()

	return etcdv3.NewClient(ctx, []string{etcdAddr}, options)
}

func GetEtcdIns(prefix string, etcdClient etcdv3.Client, logger log.Logger) *etcdv3.Instancer {
	ins, err := etcdv3.NewInstancer(etcdClient, prefix, logger)
	if err != nil {
		panic(err)
	}
	return ins
}
