package captcha

import (
	"context"
	"fmt"
	"time"

	"api-server/db/rdb"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

const (
	// 每个验证码存活5分钟
	DefaultRedisExpiration = 5 * time.Minute
	// 存在 redis 中的 key 前缀
	DefaultRedisPrefixKey = "captcha"
)

type redisStore struct {
	expiration time.Duration
	prefixKey  string
}

var store base64Captcha.Store

func GetRedisStore() base64Captcha.Store {
	if store == nil {
		store = newRedisStore(DefaultRedisExpiration, DefaultRedisPrefixKey)
	}
	return store
}

// redis store
func newRedisStore(expiration time.Duration, prefixKey string) base64Captcha.Store {
	s := new(redisStore)
	s.expiration = expiration
	s.prefixKey = prefixKey
	if s.prefixKey == "" {
		s.prefixKey = DefaultRedisPrefixKey
	}
	if s.expiration == 0 {
		s.expiration = DefaultRedisExpiration
	}
	return s
}

// set
func (s *redisStore) Set(id string, value string) error {
	c := rdb.GetClient()
	k := fmt.Sprintf("%s-%s", s.prefixKey, id)
	_, err := c.SetNX(context.Background(), k, value, s.expiration).Result()
	if err != nil {
		zap.L().Error("redis set key", zap.Error(err))
	}
	return err
}

// get
func (s *redisStore) Get(id string, clear bool) string {
	c := rdb.GetClient()
	k := fmt.Sprintf("%s-%s", s.prefixKey, id)
	v, err := c.Get(context.Background(), k).Result()
	if err != nil {
		zap.L().Error("redis get key", zap.Error(err))
		return ""
	}
	return v
}

// verify
func (s *redisStore) Verify(id, answer string, clear bool) bool {
	c := rdb.GetClient()
	k := fmt.Sprintf("%s-%s", s.prefixKey, id)
	v, err := c.Get(context.Background(), k).Result()
	if err != nil {
		zap.L().Error("redis verify key", zap.Error(err))
		return false
	}
	if v == answer {
		if clear {
			_, err = c.Del(context.Background(), k).Result()
			if err != nil {
				zap.L().Error("redis verify del key", zap.Error(err))
				return false
			}
		}
		return true
	}
	return false
}
