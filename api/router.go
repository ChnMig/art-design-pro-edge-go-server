package api

import (
	"github.com/gin-gonic/gin"

	"api-server/api/app/system/menu"
	"api-server/api/app/system/user"
	"api-server/api/middleware"
)

func systemRouter(router *gin.RouterGroup) {
	systemRouter := router.Group("/admin/system")
	{
		systemRouter.GET("/user/login/captcha", user.GetCaptcha)
		systemRouter.POST("/user/login", user.Login)
		systemRouter.GET("/user/info", middleware.TokenVerify, user.GetUserInfo)
		systemRouter.GET("/user/menu", middleware.TokenVerify, user.GetUserMenuList)
		systemRouter.GET("/menu", middleware.TokenVerify, menu.GetMenuList)
		systemRouter.POST("/menu", middleware.TokenVerify, menu.AddMenu)
		systemRouter.DELETE("/menu", middleware.TokenVerify, menu.DeleteMenu)
		systemRouter.PUT("/menu", middleware.TokenVerify, menu.UpdateMenu)
	}
}

// InitApi init gshop app
func InitApi() *gin.Engine {
	// gin.Default uses Use by default. Two global middlewares are added, Logger(), Recovery(), Logger is to print logs, Recovery is panic and returns 500
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	router.SetTrustedProxies(nil)
	// Add consent cross-domain middleware
	router.Use(middleware.CorssDomainHandler())
	// static
	router.Static("/static", "./static")
	// api-v1
	v1 := router.Group("/api/v1")
	{
		systemRouter(v1)
	}
	return router
}
