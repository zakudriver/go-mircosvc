package config

import (
	"flag"
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
)

type config struct {
	LogPath       string       `yaml:"LogPath"`
	JwtAuthSecret string       `yaml:"JwtAuthSecret"`
	PidPath       string       `yaml:"PidPath"`
	// ServerHost    string       `yaml:"ServerHost"`
	ServerAddr    string       `yaml:"ServerAddr"`
	EtcdAddr      string       `yaml:"EtcdAddr"`
	ZipkinAddr    string       `yaml:"ZipkinAddr"`
	Redis         db.RedisConf `yaml:"Redis"`
}

var (
	confFile string
	c        = new(config)
)

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
