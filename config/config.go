package config

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/Zhan9Yunhua/logger"
	"github.com/Zhan9Yunhua/blog-svr/utils"
)

type svrCfg struct {
	Env       string
	Port      int
	WsPort    int
	Cors      string
	LogDir    string
	APIPrefix string
}

type dbCfg struct {
	AuthSource string
	Username   string
	Password   string
	Host       string
	Port       int
}

type redisCfg struct {
	Password string
	Host     string
	Port     int
	PoolSize int
}

type jwtCfg struct {
	Secret    string
	OverTime  int
	CacheTime int
}

type uploadCfg struct {
	MaxSize    int
	ArticleDir string
	ProfileDir string
	Addr       string
}

type user struct {
	Salt        string
	Guest       string
	DefAvatar   string
	CodeOverdue int
}

type email struct {
	From     string
	AuthCode string
	Host     string
	Port     int
}

const defCfgFile = "./config.yml"

var (
	cfgFile string
	cfgMap  map[string]interface{}

	SvrCfg    svrCfg
	DbCfg     dbCfg
	RedisCfg  redisCfg
	JwtCfg    jwtCfg
	UploadCfg uploadCfg
	UserCfg   user
	EmailCfg  email
)

func init() {
	logger.Infoln("config init")
	flag.StringVar(&cfgFile, "cf", "", "config file path")
	flag.StringVar(&SvrCfg.Env, "env", "", "env")

	flag.Parse()
	if SvrCfg.Env == "" {
		SvrCfg.Env = "DEV"
	}

	initCfgFile()

	cfgHandler("ser", &SvrCfg)
	cfgHandler("db", &DbCfg)
	cfgHandler("redis", &RedisCfg)
	cfgHandler("jwt", &JwtCfg)
	cfgHandler("upload", &UploadCfg)
	cfgHandler("user", &UserCfg)
	cfgHandler("email", &EmailCfg)
}

func initCfgFile() {
	if cfgFile == "" {
		cfgFile = defCfgFile
	}
	path, err := filepath.Abs(cfgFile)
	if err != nil {
		logger.Fatalln(err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatalln("ReadFile:", err.Error())
	}

	if err := yaml.Unmarshal(b, &cfgMap); err != nil {
		logger.Fatalln(err)
	}
}

func cfgHandler(key string, target interface{}) {
	if err := utils.JSON2Struct(cfgMap[key].(map[interface{}]interface{}), target); err != nil {
		logger.Fatalln(err)
	}
}
