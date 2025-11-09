package api

import (
	"github.com/gin-gonic/gin"

	"api-server/api/app"
	"api-server/api/middleware"
	"api-server/config"
)

// InitApi 初始化 HTTP 服务
func InitApi() *gin.Engine {
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
