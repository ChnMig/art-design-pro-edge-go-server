package acme

import (
	"net/http"
	"testing"

	"api-server/config"
)

// TestSetup_Disabled 验证在未启用 ACME 时不会修改 server，也不会创建额外的 HTTP 服务。
func TestSetup_Disabled(t *testing.T) {
	originalEnable := config.EnableACME
	t.Cleanup(func() {
		config.EnableACME = originalEnable
	})

	config.EnableACME = false

	srv := &http.Server{}

	ctx := Setup(srv)
	if ctx == nil {
		t.Fatalf("Setup() 返回 nil")
	}
	if ctx.Enabled {
		t.Errorf("ctx.Enabled = true, want false when ACME disabled")
	}
	if ctx.HTTPServer != nil {
		t.Errorf("ctx.HTTPServer != nil, want nil when ACME disabled")
	}
	if srv.TLSConfig != nil {
		t.Errorf("srv.TLSConfig != nil, want nil when ACME disabled")
	}
}

// TestSetup_Enabled 验证在启用 ACME 时会为 server 配置 TLSConfig，并创建挑战用 HTTP 服务。
// 该测试不发起真实网络连接，仅检查配置与结构。
func TestSetup_Enabled(t *testing.T) {
	originalEnable := config.EnableACME
	originalDomain := config.ACMEDomain
	originalCacheDir := config.ACMECacheDir
	originalListenPort := config.ListenPort
	t.Cleanup(func() {
		config.EnableACME = originalEnable
		config.ACMEDomain = originalDomain
		config.ACMECacheDir = originalCacheDir
		config.ListenPort = originalListenPort
	})

	config.EnableACME = true
	config.ACMEDomain = "example.com"
	config.ACMECacheDir = t.TempDir()
	config.ListenPort = 443

	srv := &http.Server{
		Addr: ":443",
	}

	ctx := Setup(srv)
	if ctx == nil {
		t.Fatalf("Setup() 返回 nil")
	}
	if !ctx.Enabled {
		t.Fatalf("ctx.Enabled = false, want true when ACME enabled")
	}
	if ctx.HTTPServer == nil {
		t.Fatalf("ctx.HTTPServer is nil, want non-nil when ACME enabled")
	}
	if ctx.HTTPServer.Addr != ":80" {
		t.Errorf("ctx.HTTPServer.Addr = %q, want %q", ctx.HTTPServer.Addr, ":80")
	}
	if srv.TLSConfig == nil {
		t.Errorf("srv.TLSConfig is nil, want non-nil when ACME enabled")
	}
}
