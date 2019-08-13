package config

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/Zhan9Yunhua/blog-svr/utils"
	"gopkg.in/yaml.v2"

	"github.com/Zhan9Yunhua/logger"
)

const (
	defConfFile = "./gateway/config.yml"
)

type config struct {
	LogPath       string
	JwtAuthSecret string
	PidPath       string
	ServerPort    string
}

var (
	confFile string
	conf     map[string]string
	c        =new(config)
)

func init() {
	logger.Infoln("config init")
	flag.StringVar(&confFile, "cf", "", "config file path")

	flag.Parse()

	if err := initConfFile(); err != nil {
		logger.Fatalln(err)
	}

	if err := handleConf(c); err != nil {
		logger.Fatalln(err)
	}
}

func initConfFile() error {
	if confFile == "" {
		confFile = defConfFile
	}
	path, err := filepath.Abs(confFile)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, &conf)
}

func handleConf(target *config) error {
	c := map[interface{}]interface{}{}
	for k, v := range conf {
		c[k] = v
	}

	return utils.JSON2Struct(c, target)
}

func GetConfig() *config {
	return c
}
