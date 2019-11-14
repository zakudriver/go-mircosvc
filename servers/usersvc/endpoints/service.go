package endpoints

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/kum0/blog-svr/common"
	"github.com/kum0/blog-svr/servers/usersvc/endpoints/model"
	"github.com/kum0/blog-svr/shared/email"
	"github.com/kum0/blog-svr/shared/session"
	"github.com/kum0/blog-svr/shared/validator"
	"github.com/kum0/blog-svr/utils"
	"strings"
)

type IUserService interface {
	GetUser(context.Context, string) (*GetUserResponse, error)
	Login(context.Context, LoginRequest) (*LoginResponse, error)
	SendCode(context.Context) (*SendCodeResponse, error)
	Register(context.Context, RegisterRequest) error
}

func NewUserService(db *sql.DB, redis *redis.Pool, email *email.Email) IUserService {
	return &UserService{
		db,
		redis,
		email,
		session.NewSession(),
		validator.NewValidator(),
	}
}

type UserService struct {
	mysql     *sql.DB
	redis     *redis.Pool
	email     *email.Email
	session   *session.Session
	validator *validator.Validator
}

func (svc *UserService) GetUser(_ context.Context, uid string) (*GetUserResponse, error) {
	return &GetUserResponse{strings.ToUpper(uid)}, nil
}

func (svc *UserService) Login(_ context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := svc.validator.LazyValidate(req); err != nil {
		return nil, common.ArgsErr(err)
	}

	user := new(model.User)
	sql := fmt.Sprintf("SELECT `id`, `username`, `password`, `avatar`, `role_id`, `recent_time`, `created_time`, "+
		"`updated_time` "+
		"FROM `User` WHERE `username`='%s'",
		req.Username)
	err := svc.mysql.QueryRow(sql).Scan(&user.Id, &user.Username, &user.Password, &user.Avatar, &user.RoleID,
		&user.RecentTime, &user.CreatedTime, &user.UpdatedTime)
	if err != nil {
		return nil, common.ArgsErr(fmt.Sprintf("[%s]该用户名不存在", req.Username))
	}

	if user.VerifyPassword(req.Password) {
		res := new(LoginResponse)
		if err := utils.StructCopy(user, res); err != nil {
			return nil, common.ArgsErr(err)
		}
		return res, nil
	}
	return nil, common.ArgsErr("密码错误")
}

func (svc *UserService) SendCode(_ context.Context) (*SendCodeResponse, error) {
	uuid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	code := utils.NewRand(6)

	rc := svc.redis.Get()
	defer rc.Close()

	ch := make(chan error)

	go func(c chan<- error) {
		if _, err := rc.Do("SET", uuid.String(), code, "EX", 600); err != nil {
			c <- err
		}
		c <- nil
	}(ch)

	go func(c chan<- error) {
		html := fmt.Sprintf(`
      <html>
      <body>
	  <h3>
      注册码: %d
      </h3>
      </body>
      </html>
      `, code)
		c <- svc.email.Send("zy.hua1122@outlook.com", "注册码", html)
	}(ch)

	n := 2
	for c := range ch {
		n--
		if c != nil {
			close(ch)
			return nil, common.ArgsErr(c)
		}
		if n == 0 {
			close(ch)
		}
	}

	return &SendCodeResponse{uuid.String()}, nil
}

func (svc *UserService) Register(ctx context.Context, req RegisterRequest) error {
	conn := svc.redis.Get()
	defer conn.Close()

	code, err := redis.Int(conn.Do("GET", req.CodeID))
	if err != nil {
		return common.ArgsErr(err)
	}

	if code == int(req.CodeID) {
		user := new(model.User)
		pwd := user.Pwd2Md5(req.Password, user.Salt())

		sql := fmt.Sprintf("INSERT INTO `User`(`username`, `password`, `avatar`) VALUES('%s', '%s', '%s')",
			req.Username,
			pwd, "avatar")
		_, err := svc.mysql.Exec(sql)
		if err != nil {
			return common.ArgsErr(err)
		}
	} else {
		return common.ArgsErr("验证码错误")
	}

	return nil
}
