package config

import (
	"flag"
	"path/filepath"

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
	LogPath       string
	JwtAuthSecret string
	PidPath       string
	ServerAddr    string
	EtcdAddr      string
	BaseURL       string

	DBIP       string
	DBPort     int
	DBUsername string
	DBPassword string
	DBName     string
}

var (
	confFile string
	conf     map[string]interface{}
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

	if err := utils.ReadYmlFile(cf, &conf); err != nil {
		return err
	}

	if err := confMap2Stut(c); err != nil {
		return err
	}

	return nil
}

func confMap2Stut(target *config) error {
	cc := map[interface{}]interface{}{}
	for k, v := range conf {
		cc[k] = v
	}

	return utils.JSON2Struct(cc, target)
}
