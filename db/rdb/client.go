package rdb

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"api-server/config"
)

var (
	client   *redis.Client
	clientMu sync.Mutex
)

func GetClient() *redis.Client {
	db, err := getOrInitClient(context.Background())
	if err != nil {
		zap.L().Fatal("redis init failed", zap.Error(err))
	}
	return db
}

// GetInitializedClient 返回已初始化的客户端；未初始化时返回 nil，且不会触发连接。
func GetInitializedClient() *redis.Client {
	clientMu.Lock()
	defer clientMu.Unlock()
	return client
}

// Ping 使用调用方上下文检查 Redis，初始化失败时返回错误而不是终止进程。
func Ping(ctx context.Context) error {
	db, err := getOrInitClient(ctx)
	if err != nil {
		return err
	}
	return db.Ping(ctx).Err()
}

// Init 初始化 Redis 客户端连接池
func Init() error {
	_, err := getOrInitClient(context.Background())
	return err
}

func getOrInitClient(ctx context.Context) (*redis.Client, error) {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client != nil {
		return client, nil
	}
	db := redis.NewClient(&redis.Options{
		Addr:         config.RedisHost,
		Password:     config.RedisPassword,
		PoolSize:     100,
		MinIdleConns: 50,
	})
	if err := db.Ping(ctx).Err(); err != nil {
		_ = db.Close()
		zap.L().Error("redis连接失败", zap.Error(err))
		return nil, err
	}
	client = db
	return client, nil
}

func CloseClient() error {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client == nil {
		return nil
	}
	err := client.Close()
	client = nil
	return err
}
