package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysql(username, password, addr, authSource string) *sql.DB {
	host := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true", username, password, "tcp", addr, authSource)

	db, err := sql.Open("mysql", host)
	if err != nil {
		panic(fmt.Sprintf("Open mysql failed,err:%v\n", err))
	}

	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(100) // 设置最大连接数
	db.SetMaxIdleConns(16)  // 设置闲置连接数

	return db
}
