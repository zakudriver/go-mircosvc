package etcd

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
)

func NewEtcd(addr string) etcdv3.Client {
	etcdClient, err := handleEtcd(addr)
	if err != nil {
		panic(err)
	}
	return etcdClient
}

func handleEtcd(addr string) (etcdv3.Client, error) {
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

	return etcdv3.NewClient(ctx, []string{addr}, options)
}

func Register(prefix, addr string, etcdClient etcdv3.Client, logger log.Logger) *etcdv3.Registrar {
	// conf := config.GetConfig()

	// prefix := conf.Prefix         // known at compile time
	// instance := conf.ServerAddr   // taken from runtime or platform, somehow
	key := prefix + addr      // should be globally unique
	value := "http://" + addr // based on our transport

	registrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   key,
		Value: value,
	}, logger)

	registrar.Register()

	return registrar
}

func NewInstancer(addr, prefix string, logger log.Logger) *etcdv3.Instancer {
	ins, err := etcdv3.NewInstancer(NewEtcd(addr), prefix, logger)
	if err != nil {
		panic(err)
	}

	return ins
}
