package pgdb

import (
	"context"
	"sync"
	"time"

	"api-server/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	client   *gorm.DB
	clientMu sync.Mutex
)

func GetClient() *gorm.DB {
	db, err := getOrInitClient(context.Background())
	if err != nil {
		zap.L().Fatal("pgdb init failed", zap.Error(err))
	}
	return db
}

// Ping 使用调用方上下文检查 PostgreSQL，初始化失败时返回错误而不是终止进程。
func Ping(ctx context.Context) error {
	db, err := getOrInitClient(ctx)
	if err != nil {
		return err
	}
	database, err := db.DB()
	if err != nil {
		return err
	}
	return database.PingContext(ctx)
}

// CloseClient 关闭 PostgreSQL 连接池。
func CloseClient() error {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client == nil {
		return nil
	}
	database, err := client.DB()
	if err != nil {
		return err
	}
	if err := database.Close(); err != nil {
		return err
	}
	client = nil
	return nil
}

// Init 初始化 PostgreSQL 连接池。
func Init() error {
	_, err := getOrInitClient(context.Background())
	return err
}

func getOrInitClient(ctx context.Context) (*gorm.DB, error) {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client != nil {
		return client, nil
	}
	db, err := gorm.Open(postgres.Open(config.PgsqlDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用迁移时的外键约束
		CreateBatchSize:                          1000, // 批量插入大小
		DisableAutomaticPing:                     true,
	})
	if err != nil {
		zap.L().Error("连接 PostgreSQL 失败", zap.Error(err))
		return nil, err
	}
	pgDB, err := db.DB()
	if err != nil {
		zap.L().Error("获取 PostgreSQL 底层连接池失败", zap.Error(err))
		return nil, err
	}
	if err := pgDB.PingContext(ctx); err != nil {
		_ = pgDB.Close()
		return nil, err
	}
	// 设置连接池的空闲连接数、最大连接数和连接复用时长。
	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)
	client = db
	return client, nil
}
