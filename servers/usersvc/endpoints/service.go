package endpoints

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/kum0/go-mircosvc/common"
	userPb "github.com/kum0/go-mircosvc/pb/user"
	"github.com/kum0/go-mircosvc/servers/usersvc/endpoints/model"
	"github.com/kum0/go-mircosvc/shared/email"
	"github.com/kum0/go-mircosvc/shared/session"
	"github.com/kum0/go-mircosvc/shared/validator"
	"github.com/kum0/go-mircosvc/utils"
	"net/http"
	"strings"
)

type IUserService interface {
	GetUser(context.Context, string) (*userPb.GetUserResponse, error)
	Login(context.Context, LoginRequest) (*userPb.LoginResponse, error)
	SendCode(context.Context) (*userPb.SendCodeResponse, error)
	Register(context.Context, RegisterRequest) error
	UserList(context.Context, UserListRequest) (*userPb.UserListResponse, error)
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

func (svc *UserService) GetUser(_ context.Context, uid string) (*userPb.GetUserResponse, error) {
	return &userPb.GetUserResponse{Uid: strings.ToUpper(uid)}, nil
}

func (svc *UserService) Login(_ context.Context, req LoginRequest) (*userPb.LoginResponse, error) {
	if err := svc.validator.LazyValidate(req); err != nil {
		return nil, common.NewError(http.StatusBadRequest, err)
	}

	user := new(model.User)
	s := fmt.Sprintf("SELECT `id`, `username`, `password`, `avatar`, `role_id`, `recent_time`, `created_time`, "+
		"`updated_time` "+
		"FROM `User` WHERE `username`='%s'",
		req.Username)
	err := svc.mysql.QueryRow(s).Scan(&user.Id, &user.Username, &user.Password, &user.Avatar, &user.RoleID,
		&user.RecentTime, &user.CreatedTime, &user.UpdatedTime)
	if err != nil {
		return nil, common.NewError(http.StatusUnauthorized, fmt.Sprintf("[%s]该用户名不存在", req.Username))
	}

	if user.VerifyPassword(req.Password) {
		s := fmt.Sprintf("UPDATE `User` SET `recent_time` = sysdate() WHERE `id` = %d", user.Id)
		if _, err := svc.mysql.Exec(s); err != nil {
			return nil, err
		}

		res := new(userPb.LoginResponse)
		if err := utils.StructCopy(user, res); err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, common.NewError(http.StatusUnauthorized, "密码错误")
}

func (svc *UserService) SendCode(_ context.Context) (*userPb.SendCodeResponse, error) {
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

	return &userPb.SendCodeResponse{CodeID: uuid.String()}, nil
}

func (svc *UserService) Register(ctx context.Context, req RegisterRequest) error {
	conn := svc.redis.Get()
	defer conn.Close()

	code, err := redis.Int(conn.Do("GET", req.CodeID))
	if err != nil {
		return err
	}

	if code == int(req.CodeID) {
		user := new(model.User)
		pwd := user.Pwd2Md5(req.Password, user.Salt())

		sql := fmt.Sprintf("INSERT INTO `User`(`username`, `password`, `avatar`) VALUES('%s', '%s', '%s')",
			req.Username,
			pwd, "avatar")
		_, err := svc.mysql.Exec(sql)
		if err != nil {
			return err
		}
	} else {
		return common.NewError(http.StatusBadRequest, "验证码错误")
	}

	return nil
}

func (svc *UserService) UserList(ctx context.Context, req UserListRequest) (*userPb.UserListResponse, error) {
	s := fmt.Sprintf("SELECT `id`, `username`, `avatar`, `role_id`, `recent_time`, `created_time`, `updated_time` "+
		"FROM `User` WHERE id >= %d LIMIT %d;", (req.Page-1)*req.Size, req.Page)
	rs, err := svc.mysql.Query(s)
	defer rs.Close()

	if err != nil {
		return nil, err
	}

	d := make([]*userPb.UserResponse, 0)
	for rs.Next() {
		u := new(userPb.UserResponse)
		if err := rs.Scan(&u.Id, &u.Username, &u.Avatar, &u.RoleID, &u.RecentTime, &u.CreatedTime, &u.UpdatedTime);
			err != nil {
			return nil, err
		}
		d = append(d, u)
	}

	count := 0
	r := svc.mysql.QueryRow("SELECT FOUND_ROWS() AS `count`")
	if err := r.Scan(&count); err != nil {
		return nil, err
	}

	return &userPb.UserListResponse{Count: int64(count), Data: d}, nil
}
