package config

import (
	"flag"
	"github.com/Zhan9Yunhua/blog-svr/services/email"
	"path/filepath"

	"github.com/Zhan9Yunhua/blog-svr/services/db"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"github.com/Zhan9Yunhua/logger"
)

func init() {
	if err := handleConf(); err != nil {
		logger.Fatalln(err)
	}
}

const (
	DefConfFile = "servers/user/config.yml"
)

type config struct {
	LogPath       string `yaml:"LogPath"`
	JwtAuthSecret string `yaml:"JwtAuthSecret"`
	PidPath       string `yaml:"PidPath"`
	ServerAddr    string `yaml:"ServerAddr"`
	EtcdAddr      string `yaml:"EtcdAddr"`
	Prefix        string `yaml:"Prefix"`

	Mysql db.MysqlConf    `yaml:"Mysql"`
	Redis db.RedisConf    `yaml:"Redis"`
	Email email.EmailConf `yaml:"Email"`
}

var (
	confFile string
	c        = new(config)
)

func GetConfig() *config {
	return c
}

func handleConf() error {
	logger.Infoln("USER SERVER config init")
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
