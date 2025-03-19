package log

import (
	"os"
	"time"

	"api-server/config"
	runmodel "api-server/util/run-model"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// Creating Dev logger
// DEV mode outputs logs to the terminal and is more readable
func createDevLogger() *zap.Logger {
	encoder := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewTee(
		zapcore.NewSamplerWithOptions(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), os.Stdout, zap.DebugLevel), time.Second, 4, 1),
	)
	return zap.New(core, zap.AddCaller())
}

// Creating product logger
// The product pattern outputs logs to a file and is architecturally structured, in json format.
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

// SetLogger to prevent zap persistence problems after files are deleted
func SetLogger() {
	// Get log mode
	switch {
	case runmodel.IsDev():
		logger = createDevLogger()
	case runmodel.IsRelease():
		logger = createProductLogger(config.LogPath)
	default:
		logger = createProductLogger(config.LogPath)
	}
	zap.ReplaceGlobals(logger)
}

// Listen to log files
// When the log file is deleted manually, we will automatically create a new one.
func monitorFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		zap.L().Error("File listening error", zap.Error(err))
		return
	}
	defer watcher.Close()
	err = watcher.Add(config.LogPath)
	if err != nil {
		zap.L().Error("File listening error", zap.Error(err))
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Remove) {
				zap.L().Warn("the log file was deleted")
				SetLogger()
			}
			if event.Has(fsnotify.Rename) {
				zap.L().Warn("log files are renamed and new files are monitored")
				SetLogger()
			}
		case err := <-watcher.Errors:
			zap.L().Error("file listening error", zap.Error(err))
		}
	}
}

func GetLogger() *zap.Logger {
	return logger
}

func init() {
	SetLogger()
	go monitorFile()
}
