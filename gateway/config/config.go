package config

import (
	"flag"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"path/filepath"

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
	LogPath       string
	JwtAuthSecret string
	PidPath       string
	ServerPort    string
	EtcdAddr      string
}

var (
	confFile string
	conf     map[string]string
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

	if err := utils.ReadYmlFile(cf, &conf); err != nil {
		return err
	}

	if err := map2Stut(c); err != nil {
		return err
	}

	return nil
}

func map2Stut(target *config) error {
	cc := map[interface{}]interface{}{}
	for k, v := range conf {
		cc[k] = v
	}

	return utils.JSON2Struct(cc, target)
}
