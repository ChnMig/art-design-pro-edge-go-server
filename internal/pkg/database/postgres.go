package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var pgClient *gorm.DB

// InitPostgres 初始化PostgreSQL连接
func InitPostgres(dsn string) error {
	var err error
	pgClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数形式的表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// 获取底层的sql.DB对象来配置连接池
	sqlDB, err := pgClient.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("PostgreSQL connected successfully")
	return nil
}

// GetPostgres 获取PostgreSQL客户端
func GetPostgres() *gorm.DB {
	return pgClient
}

// ClosePostgres 关闭PostgreSQL连接
func ClosePostgres() error {
	if pgClient != nil {
		sqlDB, err := pgClient.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}