package log

import (
	"net/http"
	"os"
	"strings"
	"time"

	"api-server/config"
	runmodel "api-server/util/run-model"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	monitorDone chan struct{}
)

// BoundParamsKey 用于在 gin.Context 中存放已绑定的业务参数。
// 目前由 middleware.CheckParam 写入，WithRequest 读取，仅用于日志记录。
const BoundParamsKey = "__bound_params__"

func createDevLogger() *zap.Logger {
	encoder := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewTee(
		zapcore.NewSamplerWithOptions(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), os.Stdout, zap.DebugLevel), time.Second, 4, 1),
	)
	return zap.New(core, zap.AddCaller())
}

func createProductLogger(fileName string) *zap.Logger {
	fileEncoder := zap.NewProductionEncoderConfig()
	fileEncoder.EncodeTime = zapcore.ISO8601TimeEncoder
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    config.LogMaxSize,
		MaxBackups: config.LogMaxBackups,
		MaxAge:     config.LogMaxAge,
	})
	core := zapcore.NewTee(
		zapcore.NewSamplerWithOptions(
			zapcore.NewCore(zapcore.NewJSONEncoder(fileEncoder), fileWriter, zap.InfoLevel), time.Second, 4, 1),
	)
	return zap.New(core, zap.AddCaller())
}

// SetLogger 根据运行模式初始化 logger
func SetLogger() {
	switch {
	case runmodel.IsDev():
		logger = createDevLogger()
	case runmodel.IsRelease():
		logger = createProductLogger(config.LogPath)
	default:
		// 默认按开发模式处理，避免测试/包初始化阶段创建文件与目录
		logger = createDevLogger()
	}
	zap.ReplaceGlobals(logger)
}

func monitorFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		zap.L().Error("日志监听初始化失败", zap.Error(err))
		return
	}
	defer watcher.Close()
	if err = watcher.Add(config.LogPath); err != nil {
		zap.L().Error("日志监听失败", zap.Error(err))
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				zap.L().Warn("日志文件变更，重新初始化 logger")
				SetLogger()
			}
		case err := <-watcher.Errors:
			zap.L().Error("日志监听出错", zap.Error(err))
		case <-monitorDone:
			return
		}
	}
}

func GetLogger() *zap.Logger {
	if logger == nil {
		SetLogger()
	}
	return logger
}

// zapWriter 将第三方框架日志（如 gin）重定向到 zap。
type zapWriter struct {
	logger *zap.Logger
	level  zapcore.Level
}

func (w *zapWriter) Write(p []byte) (n int, err error) {
	if w == nil || w.logger == nil {
		return len(p), nil
	}
	msg := strings.TrimRight(string(p), "\r\n")
	if ce := w.logger.Check(w.level, msg); ce != nil {
		ce.Write()
	}
	return len(p), nil
}

// NewZapWriter 创建一个基于 zap 的 io.Writer，方便将第三方日志重定向到统一的 zap 日志管道。
func NewZapWriter(l *zap.Logger, level zapcore.Level) *zapWriter {
	if l == nil {
		l = GetLogger()
	}
	return &zapWriter{
		logger: l,
		level:  level,
	}
}

// StartMonitor 仅在生产环境监控日志文件
func StartMonitor() {
	if runmodel.IsRelease() {
		monitorDone = make(chan struct{})
		go monitorFile()
	}
}

// StopMonitor 退出时清理监控并刷新缓冲
func StopMonitor() {
	if monitorDone != nil {
		close(monitorDone)
	}
	if logger != nil {
		_ = logger.Sync()
	}
}

// FromContext 提供带 request 信息的 logger
func FromContext(c *gin.Context) *zap.Logger {
	if loggerVal, exists := c.Get("logger"); exists {
		if contextLogger, ok := loggerVal.(*zap.Logger); ok {
			return contextLogger
		}
	}
	return GetLogger()
}

// WithRequest 从 gin.Context 中获取带请求参数信息的 logger。
// 仅在需要排查问题时调用，避免对所有请求都记录参数。
func WithRequest(c *gin.Context) *zap.Logger {
	base := FromContext(c)

	if c == nil || c.Request == nil {
		return base
	}

	fields := []zap.Field{
		zap.String("method", c.Request.Method),
	}

	if c.Request.URL != nil {
		fields = append(fields, zap.String("path", c.Request.URL.Path))
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			fields = append(fields, zap.String("query", rawQuery))
		}
	}

	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
		if len(c.Request.PostForm) > 0 {
			fields = append(fields, zap.Any("form", c.Request.PostForm))
		}
		if c.Request.MultipartForm != nil && len(c.Request.MultipartForm.Value) > 0 {
			fields = append(fields, zap.Any("multipart_form", c.Request.MultipartForm.Value))
		}
	}

	if len(c.Params) > 0 {
		pathParams := make(map[string]string, len(c.Params))
		for _, p := range c.Params {
			pathParams[p.Key] = p.Value
		}
		fields = append(fields, zap.Any("path_params", pathParams))
	}

	if bound, exists := c.Get(BoundParamsKey); exists && bound != nil {
		fields = append(fields, zap.Any("params", bound))
	}

	return base.With(fields...)
}
