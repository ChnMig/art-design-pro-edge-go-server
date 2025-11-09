package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var v *viper.Viper

// LoadConfig 使用 Viper 加载配置
func LoadConfig() error {
	v = viper.New()

	// 设置配置文件名和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(AbsPath)
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/http-services/")

	// 支持环境变量覆盖
	v.SetEnvPrefix("HTTP_SERVICES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 默认值
	setDefaults()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zap.L().Warn("配置文件未找到，使用默认值", zap.String("path", AbsPath))
		} else {
			return fmt.Errorf("read config failed: %w", err)
		}
	} else {
		zap.L().Info("配置文件已加载", zap.String("file", v.ConfigFileUsed()))
	}

	return applyConfig()
}

func setDefaults() {
	// server
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.max_body_size", "10MB")
	v.SetDefault("server.max_header_bytes", 1<<20)
	v.SetDefault("server.shutdown_timeout", "10s")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "120s")
	v.SetDefault("server.enable_rate_limit", false)
	v.SetDefault("server.global_rate_limit", 100)
	v.SetDefault("server.global_rate_burst", 200)

	// jwt
	v.SetDefault("jwt.expiration", "12h")

	// log
	v.SetDefault("log.max_size", 50)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 30)

	// redis
	v.SetDefault("redis.host", "")
	v.SetDefault("redis.password", "")

	// postgres
	v.SetDefault("postgres.host", "")
	v.SetDefault("postgres.user", "")
	v.SetDefault("postgres.password", "")
	v.SetDefault("postgres.dbname", "")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.timezone", "Asia/Shanghai")

	// admin
	v.SetDefault("admin.password", "")
	v.SetDefault("admin.salt", "")

	// rate limit
	v.SetDefault("rate_limit.login_rate_per_minute", 5)
	v.SetDefault("rate_limit.login_burst_size", 10)
	v.SetDefault("rate_limit.general_rate_per_sec", 100)
	v.SetDefault("rate_limit.general_burst_size", 200)

	// tenant
	v.SetDefault("tenant.min_query_length", 3)
	v.SetDefault("tenant.default_code", "platform")
}

func applyConfig() error {
	// server
	ListenPort = v.GetInt("server.port")

	sizeStr := v.GetString("server.max_body_size")
	size, err := parseSize(sizeStr)
	if err != nil {
		return fmt.Errorf("invalid max_body_size: %w", err)
	}
	MaxBodySize = size

	MaxHeaderBytes = v.GetInt("server.max_header_bytes")
	ShutdownTimeout = v.GetDuration("server.shutdown_timeout")
	ReadTimeout = v.GetDuration("server.read_timeout")
	WriteTimeout = v.GetDuration("server.write_timeout")
	IdleTimeout = v.GetDuration("server.idle_timeout")
	EnableRateLimit = v.GetBool("server.enable_rate_limit")
	GlobalRateLimit = v.GetInt("server.global_rate_limit")
	GlobalRateBurst = v.GetInt("server.global_rate_burst")

	// jwt
	JWTKey = v.GetString("jwt.key")
	JWTExpiration = v.GetDuration("jwt.expiration")

	// log
	LogMaxSize = v.GetInt("log.max_size")
	LogMaxBackups = v.GetInt("log.max_backups")
	LogMaxAge = v.GetInt("log.max_age")

	// redis
	RedisHost = v.GetString("redis.host")
	RedisPassword = v.GetString("redis.password")

	// postgres -> DSN
	host := v.GetString("postgres.host")
	user := v.GetString("postgres.user")
	password := v.GetString("postgres.password")
	dbname := v.GetString("postgres.dbname")
	port := v.GetInt("postgres.port")
	sslmode := v.GetString("postgres.sslmode")
	timezone := v.GetString("postgres.timezone")

	if host != "" {
		PgsqlDSN = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			host,
			user,
			password,
			dbname,
			port,
			sslmode,
			timezone,
		)
	} else {
		PgsqlDSN = ""
	}

	// admin
	AdminPassword = v.GetString("admin.password")
	PWDSalt = v.GetString("admin.salt")

	// rate limit
	LoginRatePerMinute = v.GetInt("rate_limit.login_rate_per_minute")
	LoginBurstSize = v.GetInt("rate_limit.login_burst_size")
	GeneralRatePerSec = v.GetInt("rate_limit.general_rate_per_sec")
	GeneralBurstSize = v.GetInt("rate_limit.general_burst_size")

	// tenant
	TenantMinQueryLength = v.GetInt("tenant.min_query_length")
	DefaultTenantCode = v.GetString("tenant.default_code")
	if DefaultTenantCode == "" {
		DefaultTenantCode = "platform"
	}

	return nil
}

// WatchConfig 监听配置变化
func WatchConfig(onChange func()) {
	if v == nil {
		return
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Info("配置文件发生变化，重新加载",
			zap.String("file", e.Name),
			zap.String("op", e.Op.String()),
		)

		// 重新应用配置
		if err := applyConfig(); err != nil {
			zap.L().Error("配置热加载失败", zap.Error(err))
			return
		}

		if onChange != nil {
			onChange()
		}

		zap.L().Info("配置热加载完成")
	})
}

// GetViper 返回 viper 实例
func GetViper() *viper.Viper {
	return v
}

// parseSize 解析 KB/MB/GB 等大小字符串
func parseSize(sizeStr string) (int64, error) {
	var size int64
	var unit string
	_, err := fmt.Sscanf(sizeStr, "%d%s", &size, &unit)
	if err != nil {
		return 0, err
	}

	switch strings.ToUpper(unit) {
	case "B", "":
		return size, nil
	case "KB", "K":
		return size * 1024, nil
	case "MB", "M":
		return size * 1024 * 1024, nil
	case "GB", "G":
		return size * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unknown size unit: %s", unit)
	}
}
