package utils

import (
	"fmt"
	"time"

	// "github.com/kum0/blog-svr/db/redis"
	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Value     string
	Timestamp int64
	Claims    map[string]interface{}
}

// 生成token
func NewToken(a interface{}) *Token {

	switch t := a.(type) {
	case string:
		return &Token{
			Value:     t,
			Timestamp: time.Now().Unix(),
		}
	case map[string]interface{}:
		return &Token{
			Timestamp: time.Now().Unix(),
			Claims:    t,
		}
	default:
		return &Token{}
	}
}

// 生成token
func (t *Token) CreateValue(secret string) (string, error) {
	if t.Value != "" {
		return t.Value, nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(t.Claims))

	value, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	t.Value = value
	return value, nil
}

// 保存token对象到redis , config.JwtCfg.OverTime
// func (t *Token) Save(uid int, time int) error {
// 	rc := redis.RedisConn()
// 	defer rc.Close()
//
// 	bys, err := json.Marshal(t)
// 	if err != nil {
// 		return err
// 	}
//
// 	if _, err := rc.Do("SET", uid, bys, "EX", time); err != nil {
// 		return err
// 	}
// 	return nil
// }

// 解析token
func (t *Token) ParseToken(secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(t.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
