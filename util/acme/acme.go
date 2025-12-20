package acme

import (
	"net/http"
	"path/filepath"

	"api-server/config"
	pathtool "api-server/util/path-tool"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// Context 描述 ACME 自动 TLS（自动证书签发与续期）相关的运行时信息。
// Enabled 为 true 表示当前进程启用了 ACME；
// HTTPServer 为 ACME HTTP-01 挑战所使用的辅助 HTTP 服务（监听 80 端口），可能为 nil。
type Context struct {
	Enabled    bool
	HTTPServer *http.Server
}

// Setup 根据全局配置为主 HTTP 服务挂载 ACME 自动 TLS 能力。
// - 当未启用 ACME 时，仅返回 Disabled 的上下文，不修改传入的 server；
// - 当启用 ACME 时：
//   - 为 server 配置 TLSConfig（由 autocert.Manager 提供）；
//   - 创建 ACME 挑战用的 HTTP 服务器（监听 :80）并写入返回的 Context。
func Setup(server *http.Server) *Context {
	ctx := &Context{
		Enabled:    false,
		HTTPServer: nil,
	}

	if !config.EnableACME {
		return ctx
	}

	// 解析证书缓存目录（相对路径基于程序所在目录）
	cacheDir := config.ACMECacheDir
	if cacheDir == "" {
		cacheDir = "acme-cert-cache"
	}
	if !filepath.IsAbs(cacheDir) {
		cacheDir = filepath.Join(config.AbsPath, cacheDir)
	}
	if err := pathtool.CreateDir(cacheDir); err != nil {
		zap.L().Fatal("创建 ACME 证书缓存目录失败",
			zap.String("dir", cacheDir),
			zap.Error(err),
		)
	}

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.ACMEDomain),
		Cache:      autocert.DirCache(cacheDir),
	}

	server.TLSConfig = manager.TLSConfig()

	ctx.Enabled = true
	ctx.HTTPServer = &http.Server{
		Addr:    ":80",
		Handler: manager.HTTPHandler(nil),
	}

	zap.L().Info("ACME 自动 TLS 已启用",
		zap.String("domain", config.ACMEDomain),
		zap.Int("https_port", config.ListenPort),
		zap.String("acme_cache_dir", cacheDir),
	)

	if config.ListenPort != 443 {
		zap.L().Warn("ACME TLS 监听端口不是 443，可能导致证书签发失败（CA 通常固定连接 443）",
			zap.Int("https_port", config.ListenPort),
		)
	}

	return ctx
}
