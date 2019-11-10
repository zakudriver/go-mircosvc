package config

import (
	"github.com/kum0/blog-svr/utils"
	"strconv"
)

type config struct {
	ServiceName string `env:"SERVICE_NAME=usersvc"`
	LogPath     string `env:"LOG_PATH=./usersvc.log"`
	// HttpPort    string `env:"HTTP_PORT=5001"`
	GrpcPort    string `env:"GRPC_PORT=5002"`
	ZipkinAddr   string `env:"ZIPKIN_ADDR=http://localhost:9411"`
	// ZipkinAddr   string `env:"ZIPKIN_ADDR=http://172.20.0.1:9411/api/v2/spans"`
	RETRYMAX     string `env:"RETRY_MAX=3"`
	RetryMax     int
	RETRYTIMEOUT string `env:"RETRY_TIMEOUT=30000"`
	RetryTimeout int
	EtcdAddr     string `env:"ETCD_HOST=localhost:2379"`
	// Mysql
	MysqlUsername   string `env:"MYSQL_USERNAME=webtest"`
	MysqlPassword   string `env:"MYSQL_PASSWORD=zyhuatest"`
	MysqlAddr       string `env:"MYSQL_ADDR=118.24.103.174:3306"`
	MysqlAuthsource string `env:"MYSQL_AUTHSOURCE=webtest"`
	// Redis
	RedisAddr      string `env:"REDIS_ADDR=118.24.103.174:6300"`
	RedisPassword  string `env:"REDIS_PASSWORD=zyhua1122"`
	REDISMAXIDLE   string `env:"REDIS_MAXIDLE=30"`
	RedisMaxIdle   int
	REDISMAXACTIVE string `env:"REDIS_MAXACTIVE=30"`
	RedisMaxActive int
	// 	email
	EmailFrom     string `env:"EMAIL_FROM=zy.hua1122@qq.com"`
	EmailAuthCode string `env:"EMAIL_AUTHCODE=afyehpitqgmvbecc"`
	EmailHost     string `env:"EMAIL_HOST=smtp.qq.com"`
	EMAILPORT     string `env:"EMAIL_PORT=25"`
	EmailPort     int
	EmailSender   string `env:"EMAIL_SENDER=Kumo"`
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

	c.RetryMax = string2Int(c.RETRYMAX)
	c.RetryTimeout = string2Int(c.RETRYTIMEOUT)

	c.RedisMaxIdle = string2Int(c.REDISMAXIDLE)
	c.RedisMaxActive = string2Int(c.REDISMAXACTIVE)

	c.EmailPort = string2Int(c.EMAILPORT)
}

func string2Int(s string) int {
	r, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(err)
	}
	return int(r)
}
