package txorm

import (
	"context"
	"gitee.com/yanwenqing/backend-lib/logz"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
)

const (
	DefaultDelChan = "teco_xorm_del_cache"
)

// MemoryStore represents in-memory store
type MemoryStore struct {
	store         map[interface{}]interface{}
	printDebugLog bool
	rds           *redis.Client
	delChan       string
	mutex         sync.RWMutex
	l             *logz.Logger
}

// NewMemoryStore creates a new store in memory
func NewMemoryStore(r *redis.Client, printDebugLog bool, delChan ...string) *MemoryStore {
	l := logz.New(&logz.ServiceInfo{
		Module:         "orm.txorm.MemoryCache",
		ServiceId:      "",
		ServiceName:    "",
		ServiceVersion: "",
	})
	ret := &MemoryStore{
		store:         make(map[interface{}]interface{}),
		rds:           r,
		printDebugLog: printDebugLog,
		delChan:       DefaultDelChan,
		l:             l,
	}
	if len(delChan) > 0 {
		ret.delChan = delChan[0]
	}
	go ret.sub()
	return ret
}

// Put puts object into store
func (s *MemoryStore) Put(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.printDebugLog {
		s.l.Info("put key", zap.String("key", key), zap.Any("value", value))
	}
	s.store[key] = value
	return nil
}

// Get gets object from store
func (s *MemoryStore) Get(key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.printDebugLog {
		s.l.Info("get key", zap.String("key", key))
	}
	if v, ok := s.store[key]; ok {
		return v, nil
	}

	return nil, ErrNotExist
}

// Del deletes object
func (s *MemoryStore) Del(key string) error {
	s.rds.Publish(context.TODO(), s.delChan, key)
	return s.del(key)
}

func (s *MemoryStore) del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.printDebugLog {
		s.l.Info("del key", zap.String("key", key))
	}
	delete(s.store, key)
	return nil
}

func (s *MemoryStore) sub() {
	sub := s.rds.Subscribe(context.TODO(), s.delChan)
	for {
		msg, err := sub.ReceiveMessage(context.TODO())
		if err != nil {
			s.l.Error("redis sub error", zap.Error(err))
			continue
		}
		if s.printDebugLog {
			s.l.Info("redis sub msg", zap.String("channel", msg.Channel), zap.String("msg", msg.Payload))
		}
		_ = s.del(msg.Payload)
	}
}
