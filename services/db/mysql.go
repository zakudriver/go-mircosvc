package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Zhan9Yunhua/logger"
)

type MysqlConf struct {
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	Host       string `yaml:"Host"`
	Port       int    `yaml:"Port"`
	AuthSource string `yaml:"AuthSource"`
}

func NewMysql(conf MysqlConf) *sql.DB {
	host := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", conf.Username, conf.Password, "tcp", conf.Host,
		conf.Port, conf.AuthSource)

	db, err := sql.Open("mysql", host)
	if err != nil {
		// fmt.Printf("Open mysql failed,err:%v\n", err)
		logger.Errorf("Open mysql failed,err:%v\n", err)
		return nil
	}

	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(100) // 设置最大连接数
	db.SetMaxIdleConns(16)  // 设置闲置连接数

	return db
}
