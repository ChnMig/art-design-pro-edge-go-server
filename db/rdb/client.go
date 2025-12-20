package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"api-server/config"
)

// redisDB redis连接池
var client *redis.Client

// InitRedisClient 初始化客户端连接池
func Init() error {
	client = redis.NewClient(&redis.Options{
		Addr:         config.RedisHost,
		Password:     config.RedisPassword,
		PoolSize:     100,
		MinIdleConns: 50,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		zap.L().Error("redis连接失败", zap.Error(err))
		return err
	}
	return nil
}

func GetClient() *redis.Client {
	if client == nil {
		Init()
	}
	return client
}

func CloseClient() {
	client.Close()
	client = nil
}
