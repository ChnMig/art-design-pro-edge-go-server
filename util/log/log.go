package log

import (
	"os"
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
	default:
		logger = createProductLogger(config.LogPath)
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
