package endpoints

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
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
	return &LoginResponse{Username: req.Username, Id: 11, Avatar: "ava", RoleID: 12, RecentTime: "time"}, nil
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
			return nil, c
		}
		if n == 0 {
			close(ch)
		}
	}

	return &SendCodeResponse{uuid.String()}, nil
}
