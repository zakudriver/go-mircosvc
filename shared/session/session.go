package session

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/kum0/blog-svr/common"
	"github.com/gomodule/redigo/redis"
)

type Sessioner interface {
	Set(key string, value interface{}) // 设置Session
	Get(key string) interface{}        // 获取Session
	Del(key string)                    // 删除Session
	Sid() string                       // 当前Session ID
}

func NewSession() *Session {
	return &Session{
		MaxAge: int64(common.MaxAge),
		Data:   make(map[string]interface{}),
	}
}

type Session struct {
	SID              string                 // 唯一标示
	lock             sync.Mutex             // 一把互斥锁
	LastAccessedTime time.Time              // 最后访问时间
	MaxAge           int64                  // 超时时间
	Data             map[string]interface{} // 主数据
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

	return nil
}

func (s *Session) Del(key string) {
	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
	}
}

func (s *Session) Sid() string {
	return s.SID
}

type Storager interface {
	// SessionInit(sid string) (Session, error)
	SetSession(session *Session) error
	ReadSession(sid string) (*Session, error)
	DestroySession(sid string) error
	ExistsSession(sid string) bool
	// GCSession(maxLifeTime int64)
}

func NewStorage(pool *redis.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

type Storage struct {
	lock sync.Mutex // 一把互斥锁
	pool *redis.Pool
}

func (st *Storage) SetSession(session *Session) error {
	st.lock.Lock()
	defer st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	jsonStr, err := json.Marshal(session)
	if err != nil {
		return err
	}
	if _, err := conn.Do("SET", session.SID, string(jsonStr), string(session.MaxAge)); err != nil {
		return err
	}

	return nil
}

func (st *Storage) ReadSession(sid string) (*Session, error) {
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

func (st *Storage) DestroySession(sid string) error {
	st.lock.Lock()
	st.lock.Unlock()

	conn := st.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", sid)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) ExistsSession(sid string) bool {
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
