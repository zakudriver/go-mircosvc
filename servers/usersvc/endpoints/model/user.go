package model

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Id          int    `map:"username"`
	Username    string `map:"username"`
	Password    string
	Avatar      string    `map:"avatar"`
	RoleID      uint8     `map:"roleID"`
	RecentTime  time.Time `map:"recentTime"`
	CreatedTime time.Time `map:"createdTime"`
	UpdatedTime time.Time `map:"updatedTime"`
}

// md5加密
func (u *User) Pwd2Md5(pwd, salt string) (hash string) {
	pwd = fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	hash = salt + ":" + pwd + ":" + salt
	hash = salt + fmt.Sprintf("%x", md5.Sum([]byte(hash)))
	return
}

// 加盐
func (u *User) Salt() (salt string) {
	if u.Password != "" && len(u.Password) != 42 {
		return
	}

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
