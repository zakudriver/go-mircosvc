package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/config"
	"github.com/Zhan9Yunhua/logger"
)

var Mdb *sql.DB

func init() {
	initMysql()
}

func initMysql() {
	dbCfg := config.DbCfg
	host := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbCfg.Username, dbCfg.Password, "tcp", dbCfg.Host, dbCfg.Port, dbCfg.AuthSource)

	db, err := sql.Open("mysql", host)
	if err != nil {
		// fmt.Printf("Open mysql failed,err:%v\n", err)
		logger.Errorf("Open mysql failed,err:%v\n", err)
		return
	}

	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(100) // 设置最大连接数
	db.SetMaxIdleConns(16)  // 设置闲置连接数

	Mdb = db
}

type Sqler interface {
	Insert(stut interface{}) error
}

type db struct {
	*sql.DB
	table string
	value map[string]interface{}
}

func From(table string) Sqler {
	return &db{DB: Mdb, table: table}
}
