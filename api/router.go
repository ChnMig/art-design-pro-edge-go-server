package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"api-server/api/app"
	"api-server/api/middleware"
	"api-server/config"
	httplog "api-server/util/log"
)

// InitApi 初始化 HTTP 服务
func InitApi() *gin.Engine {
	// 将 gin 默认日志输出重定向到 zap，保持框架日志与业务日志统一
	ginLogWriter := httplog.NewZapWriter(zap.L().With(zap.String("logger", "gin")), zapcore.InfoLevel)
	ginErrorWriter := httplog.NewZapWriter(
		zap.L().With(zap.String("logger", "gin"), zap.String("stream", "stderr")),
		zapcore.ErrorLevel,
	)
	gin.DefaultWriter = ginLogWriter
	gin.DefaultErrorWriter = ginErrorWriter

	if config.RunModel == config.RunModelDevValue {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.SetTrustedProxies(nil)

	if config.EnableRateLimit {
		router.Use(middleware.IPRateLimit(config.GlobalRateLimit, config.GlobalRateBurst))
	}
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RequestID())
	router.Use(middleware.BodySizeLimit(config.MaxBodySize))
	router.Use(middleware.CorssDomainHandler())

	router.Static("/static", "./static")

	apiGroup := router.Group("/api")
	app.RegisterRoutes(apiGroup)

	return router
}
