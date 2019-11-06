package db

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func NewRedis(addr, password string, maxIdle, maxActive int) *redis.Pool {
	// addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	return &redis.Pool{
		MaxIdle:     maxIdle,            // 最大空闲连接数
		MaxActive:   maxActive,          // 最大连接数
		IdleTimeout: 1000 * time.Second, // 空闲连接超时时间
		Wait:        true,               // 如果超过最大连接，是报错，还是等待
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}

			_, err = c.Do("AUTH", password)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}
}

// func RedisConn() redis.Conn {
// 	return pool.Get()
// }
