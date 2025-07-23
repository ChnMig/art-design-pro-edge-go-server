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

	return nil
}
