package model

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"strconv"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/config"
	"github.com/Zhan9Yunhua/blog-svr/db/mysql"
)

// 权限
const (
	Root  = 0
	Guest = 1
)

type User struct {
	Id         int    `json:"id" db:"id" `
	Username   string `json:"username" db:"username"`
	Password   string `json:"password" db:"password"`
	Avatar     string `json:"avatar" db:"avatar"`
	Permission int32  `json:"permission" db:"permission"`
	CreatedAt  string `json:"createdAt" db:"createdAt"`
	UpdatedAt  string `json:"updatedAt" db:"updatedAt"`
	Token      string
}

// Password 转 md5
func (u *User) Pwd2Md5(pwd, salt string) (hash string) {
	pwd = fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	hash = pwd + ":" + config.UserCfg.Salt + ":" + salt
	hash = fmt.Sprintf("%x", md5.Sum([]byte(hash)))
	return
}

// 加盐
func (u *User) Salt() (salt string) {
	if u.Password == "" {
		salt = strconv.Itoa(int(time.Now().Unix()))
	} else {
		salt = u.Password[0:10]
	}
	return
}

// 验证密码
func (u *User) VerifyPassword(pwd string) bool {
	if pwd == "" || u.Password == "" {
		return false
	}

	return u.Pwd2Md5(pwd, "") == u.Password
}

// 注册
func (u *User) Register(username, password string, permission int32) error {
	u.Username = username
	u.Password = password
	u.Permission = permission
	u.Avatar = config.UserCfg.DefAvatar

	stmt, err := mysql.Mdb.Prepare("INSERT INTO `user`(username, password, permission) VALUES (?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(u.Username); err != nil {
		return err
	}
	return nil
}

// 生成token
func (u *User) NewToken() (string, error) {
	token := utils.NewToken(map[string]interface{}{"id": u.Id})

	tokenStr, err := token.CreateValue()

	if err != nil {
		return "", err
	}

	u.Token = tokenStr
	return tokenStr, nil
}

// 保存token到redis
func (u *User) SaveRedis() error {
	if u.Token == "" {
		return errors.New("token is empty !")
	}

	token := utils.NewToken(u.Token)
	if err := token.Save(u.Id, config.JwtCfg.OverTime); err != nil {
		return err
	}

	return nil
}
