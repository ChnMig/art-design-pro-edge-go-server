package config

import (
	"strings"

	"go.uber.org/zap"
)

const (
	minJWTKeyLength  = 32
	unsafeDefaultKey = "YOUR_SECRET_KEY_HERE"
)

// CheckConfig 校验关键配置项，缺失则 fatal 并记录日志
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
		zap.L().Fatal("JWTKey 配置缺失，请在 config.yaml 中设置")
	}
	if JWTKey == unsafeDefaultKey {
		zap.L().Fatal("JWT 密钥仍使用示例值，存在严重安全风险！请修改 config.yaml 中的 jwt.key 为强密钥")
	}
	if len(JWTKey) < minJWTKeyLength {
		zap.L().Fatal("JWT 密钥长度不足",
			zap.Int("current_length", len(JWTKey)),
			zap.Int("min_required", minJWTKeyLength),
			zap.String("suggestion", "请使用至少32字符的强密钥"),
		)
	}
	if JWTExpiration == 0 {
		zap.L().Fatal("JWTExpiration 配置缺失，请在 config.yaml 中设置 jwt.expiration")
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

	if EnableACME {
		if strings.TrimSpace(ACMEDomain) == "" {
			zap.L().Fatal("已启用 ACME，但未配置 server.acme_domain，请在 config.yaml 中设置为公网可访问的域名")
		}
	}

	if EnableTLS {
		if strings.TrimSpace(TLSCertFile) == "" || strings.TrimSpace(TLSKeyFile) == "" {
			zap.L().Fatal("已启用 TLS 证书文件模式，但未正确配置 server.tls_cert_file 或 server.tls_key_file，请在 config.yaml 中设置")
		}
	}

	if EnableACME && EnableTLS {
		zap.L().Fatal("配置错误：ACME 自动 TLS 与本地证书文件 TLS 模式不能同时启用，请二选一")
	}
}
