package txorm

import (
	"context"
	"errors"
	"gitee.com/yanwenqing/backend-lib/logz"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var (
	ErrNotExist = errors.New("record does not exist")
)

// RedisCache 是一个实现 Xorm Cacher 接口的结构体
type RedisCache struct {
	cli           *redis.Client
	ttl           time.Duration
	printDebugLog bool
	l             *logz.Logger
}

// NewRedisCache 创建一个新的 RedisCache 实例
func NewRedisCache(client *redis.Client, ttl time.Duration, printDebugLog bool) *RedisCache {
	l := logz.New(&logz.ServiceInfo{
		Module:         "orm.txorm.RedisCache",
		ServiceId:      "",
		ServiceName:    "",
		ServiceVersion: "",
	})
	return &RedisCache{
		cli:           client,
		ttl:           ttl,
		printDebugLog: printDebugLog,
		l:             l,
	}
}

// Put 放入缓存
func (c *RedisCache) Put(key string, value interface{}) error {
	if c.printDebugLog {
		c.l.Info("put key", zap.String("key", key), zap.Any("value", value))
	}
	return c.cli.Set(context.TODO(), key, value, c.ttl).Err()
}

// Get 获取缓存
func (c *RedisCache) Get(key string) (interface{}, error) {
	if c.printDebugLog {
		c.l.Info("get key", zap.String("key", key))
	}
	val, err := c.cli.Get(context.TODO(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotExist
		}
		c.l.Error("get failed", zap.Error(err), zap.String("key", key))
		return nil, err
	}
	return val, nil
}

// Del 删除缓存
func (c *RedisCache) Del(key string) error {
	if c.printDebugLog {
		c.l.Info("del key", zap.String("key", key))
	}
	return c.cli.Del(context.TODO(), key).Err()
}
