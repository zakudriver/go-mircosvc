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
	defConfFile = "./config.yml"
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
	Config   config
)

func init() {
	logger.Infoln("config init")
	flag.StringVar(&confFile, "cf", "", "config file path")

	flag.Parse()

	if err := initConfFile(); err != nil {
		logger.Fatalln(err)
	}

	if err := handleConf(&Config); err != nil {
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

func handleConf(target interface{}) error {
	 c := map[interface{}]interface{}{}
	for k, v := range conf {
		c[k] = v
	}

	return utils.JSON2Struct(c, target)
}
