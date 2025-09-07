package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
)

// YamlConfig represents the configuration structure for YAML file
type YamlConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Key        string `yaml:"key"`
		Expiration string `yaml:"expiration"`
	} `yaml:"jwt"`
	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	Postgres struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     int    `yaml:"port"`
		SSLMode  string `yaml:"sslmode"`
		TimeZone string `yaml:"timezone"`
	} `yaml:"postgres"`
	Admin struct {
		Password string `yaml:"password"`
		Salt     string `yaml:"salt"`
	} `yaml:"admin"`
	RateLimit struct {
		LoginRatePerMinute int `yaml:"login_rate_per_minute"`
		LoginBurstSize     int `yaml:"login_burst_size"`
		GeneralRatePerSec  int `yaml:"general_rate_per_sec"`
		GeneralBurstSize   int `yaml:"general_burst_size"`
	} `yaml:"rate_limit"`
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as int with fallback to default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// LoadConfig loads configuration from environment variables first, then falls back to config.yaml file
// Environment variables take precedence over YAML configuration for security
func LoadConfig() error {
	// Load from YAML file first
	configPath := filepath.Join(AbsPath, "config.yaml")
	var config YamlConfig
	
	// Try to read YAML config, but don't fail if it doesn't exist (for production environments)
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse config file: %v", err)
		}
	}

	// Apply configuration values with environment variable precedence
	
	// Server configuration
	ListenPort = getEnvInt("SERVER_PORT", config.Server.Port)
	if ListenPort == 0 {
		ListenPort = 8080 // default fallback
	}

	// JWT configuration
	JWTKey = getEnv("JWT_KEY", config.JWT.Key)
	if JWTKey == "" {
		return fmt.Errorf("JWT_KEY must be set either in environment variable or config file")
	}

	jwtExpStr := getEnv("JWT_EXPIRATION", config.JWT.Expiration)
	if jwtExpStr != "" {
		expiration, err := time.ParseDuration(jwtExpStr)
		if err != nil {
			return fmt.Errorf("invalid JWT expiration format: %v", err)
		}
		JWTExpiration = expiration
	} else {
		JWTExpiration = 12 * time.Hour // default fallback
	}

	// Redis configuration
	RedisHost = getEnv("REDIS_HOST", config.Redis.Host)
	if RedisHost == "" {
		RedisHost = "127.0.0.1:6379" // default fallback
	}
	
	RedisPassword = getEnv("REDIS_PASSWORD", config.Redis.Password)

	// PostgreSQL configuration
	pgHost := getEnv("POSTGRES_HOST", config.Postgres.Host)
	pgUser := getEnv("POSTGRES_USER", config.Postgres.User)
	pgPassword := getEnv("POSTGRES_PASSWORD", config.Postgres.Password)
	pgDBName := getEnv("POSTGRES_DBNAME", config.Postgres.DBName)
	pgPort := getEnvInt("POSTGRES_PORT", config.Postgres.Port)
	pgSSLMode := getEnv("POSTGRES_SSLMODE", config.Postgres.SSLMode)
	pgTimeZone := getEnv("POSTGRES_TIMEZONE", config.Postgres.TimeZone)

	// Set defaults if not provided
	if pgHost == "" {
		pgHost = "127.0.0.1"
	}
	if pgUser == "" {
		pgUser = "postgres"
	}
	if pgDBName == "" {
		pgDBName = "server"
	}
	if pgPort == 0 {
		pgPort = 5432
	}
	if pgSSLMode == "" {
		pgSSLMode = "disable"
	}
	if pgTimeZone == "" {
		pgTimeZone = "Asia/Shanghai"
	}

	// Validate required PostgreSQL password
	if pgPassword == "" {
		return fmt.Errorf("POSTGRES_PASSWORD must be set either in environment variable or config file")
	}

	PgsqlDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		pgHost, pgUser, pgPassword, pgDBName, pgPort, pgSSLMode, pgTimeZone)

	// Admin configuration
	AdminPassword = getEnv("ADMIN_PASSWORD", config.Admin.Password)
	if AdminPassword == "" {
		return fmt.Errorf("ADMIN_PASSWORD must be set either in environment variable or config file")
	}

	PWDSalt = getEnv("PASSWORD_SALT", config.Admin.Salt)
	if PWDSalt == "" {
		return fmt.Errorf("PASSWORD_SALT must be set either in environment variable or config file")
	}

	// Rate limiting configuration
	LoginRatePerMinute = getEnvInt("LOGIN_RATE_PER_MINUTE", config.RateLimit.LoginRatePerMinute)
	if LoginRatePerMinute == 0 {
		LoginRatePerMinute = 5 // default: 5 login attempts per minute
	}

	LoginBurstSize = getEnvInt("LOGIN_BURST_SIZE", config.RateLimit.LoginBurstSize)
	if LoginBurstSize == 0 {
		LoginBurstSize = 5 // default: allow 5 burst requests
	}

	GeneralRatePerSec = getEnvInt("GENERAL_RATE_PER_SEC", config.RateLimit.GeneralRatePerSec)
	if GeneralRatePerSec == 0 {
		GeneralRatePerSec = 10 // default: 10 requests per second
	}

	GeneralBurstSize = getEnvInt("GENERAL_BURST_SIZE", config.RateLimit.GeneralBurstSize)
	if GeneralBurstSize == 0 {
		GeneralBurstSize = 20 // default: allow 20 burst requests
	}

	return nil
}
