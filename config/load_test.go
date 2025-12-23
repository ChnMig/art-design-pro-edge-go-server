package config

import (
	"path/filepath"
	"testing"
	"time"
)

func TestParseSize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"bytes", "100B", 100, false},
		{"kilobytes", "10KB", 10 * 1024, false},
		{"megabytes", "5MB", 5 * 1024 * 1024, false},
		{"gigabytes", "2GB", 2 * 1024 * 1024 * 1024, false},
		{"lowercase kb", "10kb", 10 * 1024, false},
		{"short form k", "10K", 10 * 1024, false},
		{"short form m", "5M", 5 * 1024 * 1024, false},
		{"short form g", "2G", 2 * 1024 * 1024 * 1024, false},
		{"invalid format", "invalid", 0, true},
		{"unknown unit", "10XB", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseSize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("parseSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	tests := []struct {
		name string
		key  string
		want interface{}
	}{
		{"server port", "server.port", 8080},
		{"max body size", "server.max_body_size", "10MB"},
		{"pid file", "server.pid_file", "api-server.pid"},
		{"jwt expiration", "jwt.expiration", "12h"},
		{"log max size", "log.max_size", 50},
		{"log max backups", "log.max_backups", 3},
		{"enable rate limit", "server.enable_rate_limit", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := v.Get(tt.key)
			if got != tt.want {
				t.Fatalf("default %s = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestApplyConfig(t *testing.T) {
	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if ListenPort != 8080 {
		t.Fatalf("ListenPort = %d, want 8080", ListenPort)
	}

	if MaxBodySize != 10*1024*1024 {
		t.Fatalf("MaxBodySize = %d, want %d", MaxBodySize, 10*1024*1024)
	}

	if JWTExpiration != 12*time.Hour {
		t.Fatalf("JWTExpiration = %v, want %v", JWTExpiration, 12*time.Hour)
	}

	if filepath.Base(PidFile) != "api-server.pid" {
		t.Fatalf("PidFile = %s, want base api-server.pid", PidFile)
	}
}

func TestLoadConfigWithEnv(t *testing.T) {
	pidPath := filepath.Join(t.TempDir(), "custom.pid")
	t.Setenv("HTTP_SERVICES_SERVER_PORT", "9090")
	t.Setenv("HTTP_SERVICES_JWT_EXPIRATION", "24h")
	t.Setenv("HTTP_SERVICES_SERVER_PID_FILE", pidPath)

	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if ListenPort != 9090 {
		t.Fatalf("ListenPort = %d, want 9090 (from env)", ListenPort)
	}

	if JWTExpiration != 24*time.Hour {
		t.Fatalf("JWTExpiration = %v, want 24h (from env)", JWTExpiration)
	}

	if PidFile != pidPath {
		t.Fatalf("PidFile = %s, want %s (from env)", PidFile, pidPath)
	}
}

func TestGetViper(t *testing.T) {
	_ = LoadConfig()
	got := GetViper()
	if got == nil {
		t.Fatalf("GetViper() returned nil")
	}
	if got != v {
		t.Fatalf("GetViper() did not return the expected viper instance")
	}
}
