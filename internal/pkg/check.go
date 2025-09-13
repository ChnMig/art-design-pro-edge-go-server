package config

import (
	"go.uber.org/zap"
)

// checkConfig 校验关键配置项，缺失则 fatal 并记录日志
func CheckConfig(
	JWTKey string,
	JWTExpiration int64,
	RedisHost string,
	RedisPassword string,
	PgsqlDSN string,
	AdminPassword string,
	PWDSalt string,
) {
	if JWTKey == "" {
		zap.L().Fatal("JWTKey 配置缺失")
	}
	if JWTExpiration == 0 {
		zap.L().Fatal("JWTExpiration 配置缺失")
	}
	if RedisHost == "" {
		zap.L().Fatal("RedisHost 配置缺失")
	}
	if RedisPassword == "" {
		zap.L().Fatal("RedisPassword 配置缺失")
	}
	if PgsqlDSN == "" {
		zap.L().Fatal("PgsqlDSN 配置缺失")
	}
	if AdminPassword == "" {
		zap.L().Fatal("AdminPassword 配置缺失")
	}
	if PWDSalt == "" {
		zap.L().Fatal("PWDSalt 配置缺失")
	}
}
