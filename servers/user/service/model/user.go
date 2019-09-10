package model

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Avatar     string    `json:"avatar"`
	Permission uint8     `json:"permission"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// md5加密
func (u *User) Pwd2Md5(pwd, salt string) (hash string) {
	pwd = fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	hash = pwd + ":" + salt + ":" + salt
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

	return u.Pwd2Md5(pwd, u.Salt()) == u.Password
}
