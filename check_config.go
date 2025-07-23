package main

import (
	"go.uber.org/zap"

	"api-server/config"
)

// checkConfig 校验关键配置项，缺失则 panic 并记录日志
func checkConfig() {
	if config.JWTKey == "" {
		zap.L().Fatal("JWTKey 配置缺失")
	}
	if config.JWTExpiration == 0 {
		zap.L().Fatal("JWTExpiration 配置缺失")
	}
	if config.RedisHost == "" {
		zap.L().Fatal("RedisHost 配置缺失")
	}
	if config.RedisPassword == "" {
		zap.L().Fatal("RedisPassword 配置缺失")
	}
	if config.PgsqlDSN == "" {
		zap.L().Fatal("PgsqlDSN 配置缺失")
	}
	if config.AdminPassword == "" {
		zap.L().Fatal("AdminPassword 配置缺失")
	}
	if config.PWDSalt == "" {
		zap.L().Fatal("PWDSalt 配置缺失")
	}
}
