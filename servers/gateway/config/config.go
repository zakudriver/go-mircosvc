package config

import (
	"os"
	"strconv"
)

const (
	envSpaceName    = "SPACE_NAME"
	envServiceName  = "SERVICE_NAME"
	envLogLevel     = "LOG_LEVEL"
	envLogPath      = "LOG_PATH"
	envHTTPPort     = "HTTP_PORT"
	envGRPCPort     = "GRPC_PORT"
	envRetryMax     = "RETRY_MAX"
	envRetryTimeout = "RETRY_TIMEOUT"
	envZipkinAddr   = "ZIPKIN_ADDR"
	envEtcdHost     = "ETCD_HOST"
	envEtcdPort     = "ETCD_PORT"

	defSpaceName    = "gateway_svc"
	defServiceName  = "gateway_svc"
	defLogLevel     = "info"
	defLogPath      = "./"
	defHTTPPort     = "4001"
	defGRPCPort     = "4002"
	defClientTLS    = "false"
	defRetryMax     = "3"
	defRetryTimeout = "3000"
	defZipkinAddr   = ""
	defAppdashAddr  = ""
	defEtcdHost     = "localhost"
	defEtcdPort     = "2379"
)

type config struct {
	NameSpace    string
	ServiceName  string
	LogLevel     string
	LogPath      string
	HttpPort     string
	GrpcPort     string
	ZipkinAddr   string
	RetryMax     int
	RetryTimeout int
	RouterMap    map[string]string
	EtcdHost     string
	EtcdPort     string
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
	retry, err := strconv.ParseInt(handleEnv(envRetryMax, defRetryMax), 10, 0)
	if err != nil {
		panic("RetryMax: " + err.Error())
	}

	retryTimeout, err := strconv.ParseInt(handleEnv(envRetryTimeout, defRetryTimeout), 10, 0)
	if err != nil {
		panic("RetryTimeout: " + err.Error())
	}

	c.NameSpace = handleEnv(envSpaceName, defSpaceName)
	c.ServiceName = handleEnv(envServiceName, defServiceName)
	c.LogLevel = handleEnv(envLogLevel, defLogLevel)
	c.LogPath = handleEnv(envLogPath, defLogPath)
	c.HttpPort = handleEnv(envHTTPPort, defHTTPPort)
	c.GrpcPort = handleEnv(envGRPCPort, defGRPCPort)
	c.ZipkinAddr = handleEnv(envZipkinAddr, defZipkinAddr)
	c.RetryMax = int(retry)
	c.RetryTimeout = int(retryTimeout)
	c.EtcdHost = handleEnv(envEtcdHost, defEtcdHost)
	c.EtcdPort = handleEnv(envEtcdPort, defEtcdPort)
}

func handleEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
