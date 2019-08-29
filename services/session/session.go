package session

import (
	"encoding/json"
	"errors"
	"github.com/Zhan9Yunhua/blog-svr/common"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

type Sessioner interface {
	Set(key string, value interface{}) // 设置Session
	Get(key string) interface{}        // 获取Session
	Del(key string) error              // 删除Session
	GetName() string                   // 当前Session ID
}

func NewSession() *Session {
	return &Session{
		MaxAge: int64(common.MaxAge),
		Data:   make(map[string]interface{}),
	}
}

type Session struct {
	Name             string                 // 唯一标示
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

func (s *Session) Del(key string) error {
	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
		return nil
	}
	return errors.New("key is not exist")
}

func (s *Session) GetName() string {
	return s.Name
}

type Storager interface {
	// SessionInit(sid string) (Session, error)
	SetSession(session Session) error
	ReadSession(name string) (Session, error)
	DestroySession(name string) error
	GCSession(maxLifeTime int64)
}

func NewStorager(pool *redis.Pool) *Storage {
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
	if _, err := conn.Do("SET", session.Name, string(jsonStr), string(session.MaxAge)); err != nil {
		return err
	}

	return nil
}
