package service

import (
	"database/sql"
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/services/email"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"github.com/gomodule/redigo/redis"
	"strings"
)

type UserServicer interface {
	Login(loginRequest) (string, error)
	GetUser(string) (string, error)
}

func NewUserService(mdb *sql.DB, rd *redis.Pool, email *email.Email) *UserService {
	return &UserService{
		mdb,
		rd,
		email,
	}
}

type UserService struct {
	mdb   *sql.DB
	rd    *redis.Pool
	email *email.Email
}

func (u *UserService) GetUser(s string) (string, error) {
	if s == "" {
		return "", common.ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (u *UserService) Login(params loginRequest) (string, error) {
	return "", nil
}

func (u *UserService) Register(params registerRequest) {

}

func (u *UserService) SendEmail() error {
	uuid, err := utils.NewUUID()
	if err != nil {
		return nil
	}

	code := utils.NewRand(6)

	rc := u.rd.Get()
	defer rc.Close()

	ch := make(chan error)

	go func() {
		if _, err := rc.Do("SET", uuid, code, "EX", 600); err != nil {
			ch <- err
		}
		ch <- nil
	}()

	html := fmt.Sprintf(`
      <html>
      <body>
	  <h3>
      注册码: %d
      </h3>
      </body>
      </html>
      `, code)
	go func() {
		err := u.email.Send("zy.hua1122@outlook.com", "注册码", html)
		if err != nil {
			ch <- err
		}
		ch <- nil
	}()

	return <-ch
}
