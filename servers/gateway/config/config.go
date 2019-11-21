package config

import (
	"github.com/kum0/go-mircosvc/utils"
)

type config struct {
	ServiceName  string `env:"SERVICE_NAME=gateway-svc"`
	LogPath      string `env:"LOG_PATH=./log/gateway.log"`
	HttpPort     string `env:"HTTP_PORT=4001"`
	GrpcPort     string `env:"GRPC_PORT=4002"`
	ZipkinAddr   string `env:"ZIPKIN_ADDR=http://localhost:9411/api/v2/spans"`
	RETRYMAX     string `env:"RETRY_MAX=3"`
	RETRYTIMEOUT string `env:"RETRY_TIMEOUT=30000"`
	EtcdAddr     string `env:"ETCD_HOST=localhost:2379"`
	RetryMax     int
	RetryTimeout int
	// Redis
	RedisAddr      string `env:"REDIS_ADDR=118.24.103.174:6300"`
	RedisPassword  string `env:"REDIS_PASSWORD=zyhua1122"`
	REDISMAXIDLE   string `env:"REDIS_MAXIDLE=30"`
	RedisMaxIdle   int
	REDISMAXACTIVE string `env:"REDIS_MAXACTIVE=30"`
	RedisMaxActive int

	Origin string `env:"ANY_ORIGIN=*"`
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

	c.RetryMax = utils.String2Int(c.RETRYMAX)
	c.RetryTimeout = utils.String2Int(c.RETRYTIMEOUT)

	c.RedisMaxIdle = utils.String2Int(c.REDISMAXIDLE)
	c.RedisMaxActive = utils.String2Int(c.REDISMAXACTIVE)
}
