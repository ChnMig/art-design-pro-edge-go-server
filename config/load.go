package config

import (
	"fmt"
	"os"
	"path/filepath"
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
    Tenant struct {
        MinQueryLength int    `yaml:"min_query_length"`
        DefaultCode    string `yaml:"default_code"`
    } `yaml:"tenant"`
}

// LoadConfig loads configuration from config.yaml file
func LoadConfig() error {
	configPath := filepath.Join(AbsPath, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	var config YamlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}
	// Apply configuration values
	if config.Server.Port != 0 {
		ListenPort = config.Server.Port
	}
	if config.JWT.Key != "" {
		JWTKey = config.JWT.Key
	}
	if config.JWT.Expiration != "" {
		expiration, err := time.ParseDuration(config.JWT.Expiration)
		if err != nil {
			return fmt.Errorf("invalid JWT expiration format: %v", err)
		}
		JWTExpiration = expiration
	}
	if config.Redis.Host != "" {
		RedisHost = config.Redis.Host
	}
	if config.Redis.Password != "" {
		RedisPassword = config.Redis.Password
	}
	if config.Postgres.Host != "" {
		PgsqlDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			config.Postgres.Host,
			config.Postgres.User,
			config.Postgres.Password,
			config.Postgres.DBName,
			config.Postgres.Port,
			config.Postgres.SSLMode,
			config.Postgres.TimeZone,
		)
	}
	if config.Admin.Password != "" {
		AdminPassword = config.Admin.Password
	}
    if config.Admin.Salt != "" {
        PWDSalt = config.Admin.Salt
    }
    
    // Rate limiting configuration with defaults
	if config.RateLimit.LoginRatePerMinute != 0 {
		LoginRatePerMinute = config.RateLimit.LoginRatePerMinute
	} else {
		LoginRatePerMinute = 5 // default: 5 login attempts per minute
	}
	
	if config.RateLimit.LoginBurstSize != 0 {
		LoginBurstSize = config.RateLimit.LoginBurstSize
	} else {
		LoginBurstSize = 10 // default: allow 10 burst requests
	}
	
	if config.RateLimit.GeneralRatePerSec != 0 {
		GeneralRatePerSec = config.RateLimit.GeneralRatePerSec
	} else {
		GeneralRatePerSec = 100 // default: 100 requests per second
	}
	
	if config.RateLimit.GeneralBurstSize != 0 {
		GeneralBurstSize = config.RateLimit.GeneralBurstSize
	} else {
		GeneralBurstSize = 200 // default: allow 200 burst requests
    }
    
    // Tenant configuration with defaults
    if config.Tenant.MinQueryLength != 0 {
        TenantMinQueryLength = config.Tenant.MinQueryLength
    } else {
        TenantMinQueryLength = 3
    }
    if config.Tenant.DefaultCode != "" {
        DefaultTenantCode = config.Tenant.DefaultCode
    } else {
        DefaultTenantCode = "platform"
    }
    
    return nil
}
