package config

import (
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"strconv"
)

type config struct {
	ServiceName  string `env:"SERVICE_NAME=usersvc"`
	LogPath      string `env:"LOG_PATH=./usersvc.log"`
	HttpPort     string `env:"HTTP_PORT=5001"`
	GrpcPort     string `env:"GRPC_PORT=5002"`
	ZipkinAddr   string `env:"ZIPKIN_ADDR=http://localhost:9411/api/v2/spans"`
	RETRYMAX     string `env:"RETRY_MAX=3"`
	RETRYTIMEOUT string `env:"RETRY_TIMEOUT=30000"`
	EtcdHost     string `env:"ETCD_HOST=localhost"`
	EtcdPort     string `env:"ETCD_PORT=2379"`
	RetryMax     int
	RetryTimeout int
}

var c *config

func init() {
	initConfig()
}

func GetConfig() *config {
	return c
}

func initConfig() {
	c = new(config)

	if err := utils.ParseEnvForTag(c, "env"); err != nil {
		panic(err)
	}

	retryMax, err := strconv.ParseInt(c.RETRYMAX, 10, 0)
	if err != nil {
		panic(err)
	}
	c.RetryMax = int(retryMax)

	retryTimeout, err := strconv.ParseInt(c.RETRYTIMEOUT, 10, 0)
	if err != nil {
		panic(err)
	}
	c.RetryTimeout = int(retryTimeout)
}
