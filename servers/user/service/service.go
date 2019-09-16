package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/servers/user/service/model"
	"github.com/Zhan9Yunhua/blog-svr/services/session"

	"github.com/Zhan9Yunhua/blog-svr/services/validator"

	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/Zhan9Yunhua/blog-svr/services/email"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"github.com/gomodule/redigo/redis"
)

type IUserService interface {
	Login(loginRequest) (common.ResponseData, error)
	SendCode() (common.ResponseData, error)
	Register(registerRequest) error
	Validate(interface{}) error
	GetUser(string) (string, error)
	GetUserList() (common.ResponseData, error)
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

// 登录
func (u *UserService) Login(params loginRequest) (common.ResponseData, error) {
	user := new(model.User)
	sql := fmt.Sprintf("SELECT `id`, `username`, `password`, `avatar`, `role_id`, `recent_time`, `created_time`, "+
		"`updated_time` "+
		"FROM `User` WHERE `username`='%s'",
		params.Username)
	err := u.mdb.QueryRow(sql).Scan(&user.Id, &user.Username, &user.Password, &user.Avatar, &user.RoleID,
		&user.RecentTime, &user.CreatedTime, &user.UpdatedTime)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]该用户名不存在", params.Username))
	}

	if user.VerifyPassword(params.Password) {
		return utils.Struct2MapFromTag(user), nil
	}
	return nil, errors.New("密码错误")
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
		sql := fmt.Sprintf("INSERT INTO `User`(`username`, `password`, `avatar`) VALUES('%s', '%s', '%s')",
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

	return map[string]interface{}{"codeID": uuid.String()}, nil
}

func (u *UserService) Validate(a interface{}) error {
	return u.validator.LazyValidate(a)
}

func (u *UserService) GetUser(s string) (string, error) {
	if s == "" {
		return "", common.ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (u *UserService) GetUserList() (common.ResponseData, error) {
	row, err := u.mdb.Query("SELECT `username`, `avatar`, `role`, `recent_time` FROM `User`")
	defer row.Close()
	if err != nil {
		return nil, err
	}

	r := make([]map[string]interface{}, 0)
	for row.Next() {
		var username string
		var avatar string
		var role uint
		var recentTime time.Time

		if err := row.Scan(&username, &avatar, &role); err != nil {
			return nil, err
		}
		m := map[string]interface{}{
			"username":   username,
			"avatar":     avatar,
			"role":       role,
			"recentTime": recentTime,
		}

		r = append(r, m)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
