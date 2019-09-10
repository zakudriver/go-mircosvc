package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/service/model"
	"github.com/Zhan9Yunhua/blog-svr/services/session"

	"github.com/Zhan9Yunhua/blog-svr/services/validator"

	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/services/email"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"github.com/gomodule/redigo/redis"
)

type UserServicer interface {
	Login(loginRequest) (string, error)
	GetUser(string) (string, error)
	SendCode() (common.ResponseData, error)
	Register(registerRequest) error
	Validate(interface{}) error
}

func NewUserService(mdb *sql.DB, rd *redis.Pool, email *email.Email) *UserService {
	return &UserService{
		mdb,
		rd,
		email,
		session.NewSession(),
		validator.NewValidator(),
	}
}

type UserService struct {
	mdb       *sql.DB
	rd        *redis.Pool
	email     *email.Email
	session   *session.Session
	validator *validator.Validator
}

func (u *UserService) GetUser(s string) (string, error) {
	if s == "" {
		return "", common.ErrEmpty
	}
	return strings.ToUpper(s), nil
}

// 登录
func (u *UserService) Login(params loginRequest) (*model.User, error) {
	sql := fmt.Sprintf("SELECT * FROM `user` WHERE username=%s", params.Username)
	_ := u.mdb.QueryRow(sql)

	user := new(model.User)
	user.Username = params.Username
	user.Password = params.Password

	return nil, nil
}

// 注册
func (u *UserService) Register(params registerRequest) error {
	conn := u.rd.Get()
	defer conn.Close()

	code, err := redis.Int(conn.Do("GET", params.CodeID))
	if err != nil {
		return nil
	}

	user := new(model.User)
	pwd := user.Pwd2Md5(params.Password, user.Salt())

	if code == params.Code {
		sql := fmt.Sprintf("INSERT INTO `user`(`username`, `password`, `avatar`) VALUES('%s', '%s', '%s')",
			params.Username,
			pwd, "avatar")
		_, err := u.mdb.Exec(sql)
		if err != nil {
			return err
		}
	} else {
		return errors.New("验证码错误")
	}

	return nil
}

func (u *UserService) SendCode() (common.ResponseData, error) {
	uuid, err := utils.NewUUID()
	if err != nil {
		return nil, nil
	}

	code := utils.NewRand(6)

	rc := u.rd.Get()
	defer rc.Close()

	ch := make(chan error)

	go func(c chan<- error) {
		if _, err := rc.Do("SET", uuid.String(), code, "EX", 600); err != nil {
			c <- err
		}
		c <- nil
	}(ch)

	html := fmt.Sprintf(`
      <html>
      <body>
	  <h3>
      注册码: %d
      </h3>
      </body>
      </html>
      `, code)
	go func(c chan<- error) {
		c <- u.email.Send("zy.hua1122@outlook.com", "注册码", html)
	}(ch)

	n := 2
	for c := range ch {
		n--
		if c != nil {
			close(ch)
			return nil, c
		}
		if n == 0 {
			close(ch)
		}
	}

	return common.ResponseData{"codeID": uuid.String()}, nil
}

func (u *UserService) Validate(a interface{}) error {
	return u.validator.LazyValidate(a)
}
