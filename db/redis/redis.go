package redis

import (
	"fmt"

	"github.com/Zhan9Yunhua/logger"
	"github.com/Zhan9Yunhua/blog-svr/config"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func init() {
	initRedis()
}

func initRedis() {
	addr := fmt.Sprintf("%s:%d", config.RedisCfg.Host, config.RedisCfg.Port)
	opt := redis.DialPassword(config.RedisCfg.Password)

	pool = redis.NewPool(func() (conn redis.Conn, err error) {
		conn, err = redis.Dial("tcp", addr, opt)
		if err != nil {
			logger.Errorln("Connect to redis error", err)
			return
		}
		return conn, err
	}, config.RedisCfg.PoolSize)
}

func RedisConn() redis.Conn {
	return pool.Get()
}
