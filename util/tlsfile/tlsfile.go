package tlsfile

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"sync/atomic"
	"time"

	"api-server/config"
	pathtool "api-server/util/path-tool"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// Context 描述基于本地证书文件的 TLS 运行时信息。
// Enabled 为 true 表示当前进程启用了证书文件 TLS 模式。
type Context struct {
	Enabled bool
}

var currentCert atomic.Value // *tls.Certificate

// loadCertificate 从指定路径加载证书与私钥，并更新全局证书指针。
func loadCertificate(certPath, keyPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}
	currentCert.Store(&cert)
	return nil
}

// getCurrentCertificate 返回当前生效的 TLS 证书；仅用于内部与测试。
func getCurrentCertificate() *tls.Certificate {
	value := currentCert.Load()
	if value == nil {
		return nil
	}
	cert, ok := value.(*tls.Certificate)
	if !ok {
		return nil
	}
	return cert
}

// Setup 根据全局配置为 HTTP 服务器挂载基于本地证书文件的 TLS 能力，并启动文件变更监听实现证书热更新。
// - 当未启用 TLS 证书文件模式时，仅返回 Disabled 的上下文，不修改传入的 server；
// - 当启用时：
//   - 解析证书与私钥路径（支持相对路径，相对 config.AbsPath）；
//   - 加载证书写入全局缓存；
//   - 设置 server.TLSConfig.GetCertificate 回调；
//   - 使用 fsnotify 监听证书与私钥文件变更，变更时自动重新加载。
func Setup(server *http.Server) *Context {
	ctx := &Context{
		Enabled: false,
	}

	if !config.EnableTLS {
		return ctx
	}

	certPath := config.TLSCertFile
	keyPath := config.TLSKeyFile

	if certPath == "" || keyPath == "" {
		zap.L().Fatal("已启用 TLS 证书文件模式，但未配置证书或私钥路径",
			zap.String("tls_cert_file", certPath),
			zap.String("tls_key_file", keyPath),
		)
	}

	if !filepath.IsAbs(certPath) {
		certPath = filepath.Join(config.AbsPath, certPath)
	}
	if !filepath.IsAbs(keyPath) {
		keyPath = filepath.Join(config.AbsPath, keyPath)
	}

	_ = pathtool.CreateDir(filepath.Dir(certPath))
	_ = pathtool.CreateDir(filepath.Dir(keyPath))

	if err := loadCertificate(certPath, keyPath); err != nil {
		zap.L().Fatal("加载 TLS 证书失败",
			zap.String("cert_file", certPath),
			zap.String("key_file", keyPath),
			zap.Error(err),
		)
	}

	server.TLSConfig = &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert := getCurrentCertificate()
			if cert == nil {
				return nil, fmt.Errorf("no TLS certificate loaded")
			}
			return cert, nil
		},
	}

	ctx.Enabled = true

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		zap.L().Error("创建 TLS 证书文件监听失败，后续证书将无法自动热更新",
			zap.Error(err),
		)
		return ctx
	}

	paths := []string{certPath, keyPath}
	for _, p := range paths {
		if err := watcher.Add(p); err != nil {
			zap.L().Error("监听 TLS 证书文件失败",
				zap.String("path", p),
				zap.Error(err),
			)
		}
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
					time.Sleep(200 * time.Millisecond)
					if err := loadCertificate(certPath, keyPath); err != nil {
						zap.L().Error("重新加载 TLS 证书失败",
							zap.String("cert_file", certPath),
							zap.String("key_file", keyPath),
							zap.Error(err),
						)
					} else {
						zap.L().Info("TLS 证书已重新加载",
							zap.String("cert_file", certPath),
							zap.String("key_file", keyPath),
						)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				zap.L().Error("TLS 证书文件监听错误", zap.Error(err))
			}
		}
	}()

	return ctx
}
