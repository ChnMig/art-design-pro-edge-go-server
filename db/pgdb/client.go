package pgdb

import (
	"time"

	"api-server/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var client *gorm.DB

func GetClient() *gorm.DB {
	if client == nil {
		Init()
	}
	return client
}

// Connect to the database
func Init() error {
	db, err := gorm.Open(postgres.Open(config.PgsqlDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用迁移时的外键约束
		CreateBatchSize:                          1000, // 批量插入大小
	})
	if err != nil {
		zap.L().Error("connect to mysql failed", zap.Error(err))
		return err
	}
	pgDB, err := db.DB()
	if err != nil {
		zap.L().Error("get db failed", zap.Error(err))
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	pgDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	pgDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	pgDB.SetConnMaxLifetime(time.Hour)
	client = db
	return nil
}
