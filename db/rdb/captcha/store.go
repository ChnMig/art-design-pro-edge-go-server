package captcha

import (
	"context"
	"fmt"
	"sync"
	"time"

	"api-server/db/rdb"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
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

var verifyScript = redis.NewScript(`
local value = redis.call("GET", KEYS[1])
if not value then
  return 0
end
if value ~= ARGV[1] then
  return 0
end
if ARGV[2] == "1" then
  redis.call("DEL", KEYS[1])
end
return 1
`)

var getScript = redis.NewScript(`
local value = redis.call("GET", KEYS[1])
if not value then
  return false
end
if ARGV[1] == "1" then
  redis.call("DEL", KEYS[1])
end
return value
`)

var (
	store     base64Captcha.Store
	storeOnce sync.Once
)

func GetRedisStore() base64Captcha.Store {
	storeOnce.Do(func() {
		store = newRedisStore(DefaultRedisExpiration, DefaultRedisPrefixKey)
	})
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
	created, err := c.SetNX(context.Background(), k, value, s.expiration).Result()
	if err != nil {
		zap.L().Error("redis set key", zap.Error(err))
		return err
	}
	if !created {
		return fmt.Errorf("captcha id already exists")
	}
	return nil
}

// get
func (s *redisStore) Get(id string, clear bool) string {
	c := rdb.GetClient()
	k := fmt.Sprintf("%s-%s", s.prefixKey, id)
	clearValue := "0"
	if clear {
		clearValue = "1"
	}
	v, err := getScript.Run(context.Background(), c, []string{k}, clearValue).Text()
	if err != nil {
		if err != redis.Nil {
			zap.L().Error("redis get key", zap.Error(err))
		}
		return ""
	}
	return v
}

// verify
func (s *redisStore) Verify(id, answer string, clear bool) bool {
	c := rdb.GetClient()
	k := fmt.Sprintf("%s-%s", s.prefixKey, id)
	clearValue := "0"
	if clear {
		clearValue = "1"
	}
	verified, err := verifyScript.Run(context.Background(), c, []string{k}, answer, clearValue).Int()
	if err != nil {
		zap.L().Error("redis verify key", zap.Error(err))
		return false
	}
	return verified == 1
}
