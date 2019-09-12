package model

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Id          int
	Username    string `json:"username"`
	Password    string
	Avatar      string    `json:"avatar"`
	RoleID      uint8     `json:"roleID"`
	RecentTime  time.Time `json:"recentTime"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
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
