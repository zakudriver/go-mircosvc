package config

import "os"

const (
	envLogLevel       = "GATEWAY_LOG_LEVEL"
	envHTTPPort       = "GATEWAY_HTTP_PORT"
	envGRPCPort       = "GATEWAY_GRPC_PORT"
	envClientTLS      = "GATEWAY_CLIENT_TLS"
	envServerCert     = "GATEWAY_SERVER_CERT"
	envServerKey      = "GATEWAY_SERVER_KEY"
	envRetryMax       = "GATEWAY_RETRY_MAX"
	envRetryTimeout   = "GATEWAY_RETRY_TIMEOUT"
	envZipkinV1URL    = "ZIPKIN_V1_URL"
	envZipkinV2URL    = "ZIPKIN_V2_URL"
	envLightstepToken = "GATEWAY_LIGHT_STEP_TOKEN"
	envAppdashAddr    = "GATEWAY_APPDASH_ADDR"
	envConsulHost     = "CONSULT_HOST"
	envconsultPort    = "CONSULT_PORT"
	envEtcdHost       = "ETCD_HOST"
	envEtcdPort       = "ETCD_PORT"
)

type Config struct {
	NameSpace    string
	ServiceName  string
	LogLevel     string
	LogPath      string
	ServiceHost  string
	HttpPort     string
	GrpcPort     string
	ZipkinURL    string
	RetryMax     int64
	RetryTimeout int64
	RouterMap    map[string]string
}

func init() {

}

func handleConfig() {

}

func handleEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
