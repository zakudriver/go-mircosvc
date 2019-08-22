package etcd

import (
	"context"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/Zhan9Yunhua/logger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
)

func NewEtcd() etcdv3.Client {
	etcdClient, err := handleEtcd()
	if err != nil {
		logger.Fatalln(err)
	}
	return etcdClient
}

func handleEtcd() (etcdv3.Client, error) {
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

func Register(etcdClient etcdv3.Client, logger log.Logger) *etcdv3.Registrar {
	conf := config.GetConfig()

	prefix := "/svc/user/"        // known at compile time
	instance := conf.ServerAddr     // taken from runtime or platform, somehow
	key := prefix + instance      // should be globally unique
	value := "http://" + instance // based on our transport

	registrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   key,
		Value: value,
	}, logger)

	registrar.Register()

	return registrar
}
