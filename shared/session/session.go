package session

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Sessioner interface {
	Set(key string, value interface{}) // 设置Session值
	Get(key string) interface{}        // 获取Session值
	Del(key string)                    // 删除Session值
}

type Session struct {
	SID          string
	CookieName   string
	lock         sync.Mutex
	AccessedTime time.Time // 最后访问时间
	MaxAge       int       // 超时时间
	Data         map[string]interface{}
}

func (s *Session) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Data[key] = value
}

func (s *Session) Get(key string) interface{} {
	if v, ok := s.Data[key]; ok {
		return v
	}

	return ""
}

func (s *Session) Del(key string) {
	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
	}
}

type Storager interface {
	NewCookie(session *Session) *http.Cookie
	NewSession(sid, cookieName string, maxAge int) *Session
	Save(session *Session) error
	Read(sid string) (*Session, error)
	Destroy(sid string) error
	Exists(sid string) bool
	Update(sid, t string) error
}

func NewStorage(pool *redis.Pool) Storager {
	return &Storage{
		pool: pool,
	}
}

type Storage struct {
	lock sync.Mutex
	pool *redis.Pool
}

func (st *Storage) NewSession(sid, cookieName string, maxAge int) *Session {
	return &Session{
		SID:          sid,
		CookieName:   cookieName,
		MaxAge:       maxAge,
		AccessedTime: time.Now(),
		Data:         make(map[string]interface{}),
	}
}

func (st *Storage) NewCookie(session *Session) *http.Cookie {
	return &http.Cookie{
		Name:     session.CookieName,
		Value:    session.SID,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   session.MaxAge,
	}
}

func (st *Storage) Save(session *Session) error {
	st.lock.Lock()
	defer st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	jsonStr, err := json.Marshal(session)
	if err != nil {
		return err
	}
	if _, err := conn.Do("SET", session.SID, string(jsonStr), "EX", strconv.Itoa(session.MaxAge)); err != nil {
		return err
	}

	return nil
}

func (st *Storage) Read(sid string) (*Session, error) {
	st.lock.Lock()
	defer st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	r, err := redis.Bytes(conn.Do("GET", sid))
	if err != nil {
		return nil, err
	}

	se := new(Session)
	if err := json.Unmarshal(r, se); err != nil {
		return nil, err
	}

	return se, nil
}

func (st *Storage) Destroy(sid string) error {
	st.lock.Lock()
	st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", sid)
	return err
}

func (st *Storage) Exists(sid string) bool {
	st.lock.Lock()
	defer st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	is, err := redis.Bool(conn.Do("GET", sid))
	if err != nil {
		return false
	}

	return is
}

func (st *Storage) Update(sid string, t string) error {
	st.lock.Lock()
	defer st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", sid, t)
	return err
}
