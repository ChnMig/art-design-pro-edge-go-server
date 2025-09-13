package database

import (
	"log"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(host, password string, db int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	log.Println("Redis connected successfully")
	return nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}