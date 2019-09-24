package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Zhan9Yunhua/blog-svr/shared/db"
	"github.com/Zhan9Yunhua/blog-svr/utils"

	"github.com/Zhan9Yunhua/logger"
)

func init() {
	if err := handleConf(); err != nil {
		logger.Fatalln(err)
	}
}

const (
	DefConfFile = "./gateway/config.yml"

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

type config struct {
	LogPath       string `yaml:"LogPath"`
	JwtAuthSecret string `yaml:"JwtAuthSecret"`
	PidPath       string `yaml:"PidPath"`
	// ServerHost    string       `yaml:"ServerHost"`
	ServerAddr string       `yaml:"ServerAddr"`
	EtcdAddr   string       `yaml:"EtcdAddr"`
	ZipkinAddr string       `yaml:"ZipkinAddr"`
	Redis      db.RedisConf `yaml:"Redis"`
}

type Config struct {
	LogLevel       string
	ClientTLS      bool
	ConsulHost     string
	ConsultPort    string
	HttpPort       string
	GrpcPort       string
	ServerCert     string
	ServerKey      string
	RetryMax       int64
	RetryTimeout   int64
	ZipkinV1URL    string
	ZipkinV2URL    string
	LightstepToken string
	AppdashAddr    string

	Redis db.RedisConf
}

var (
	confFile string
	c        = new(config)
)

func handleEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func GetConfig() *config {
	return c
}

func handleConf() error {
	logger.Infoln("GATEWAY config init")
	flag.StringVar(&confFile, "cf", "", "config file path")

	flag.Parse()

	cf := DefConfFile
	if confFile != "" {
		cf = confFile
	}

	cf, err := filepath.Abs(cf)
	if err != nil {
		return err
	}

	if err := utils.ReadYmlFile(cf, &c); err != nil {
		return err
	}

	return nil
}
