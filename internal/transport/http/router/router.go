package router

import (
	"github.com/gin-gonic/gin"

	systemHandler "api-server/internal/handler/system"
	systemService "api-server/internal/service/system"
	"api-server/internal/transport/http/middleware"
)

// systemRouter 配置系统路由，使用依赖注入的Handler
func systemRouter(router *gin.RouterGroup) {
	// 创建Service层
	userService := systemService.NewUserService()
	tenantService := systemService.NewTenantService()

	// 创建Handler层，注入Service依赖
	userHandler := systemHandler.NewUserHandler(userService)
	tenantHandler := systemHandler.NewTenantHandler(tenantService)

	systemRouter := router.Group("/admin/system")
	{
		// 用户认证相关路由
		systemRouter.GET("/user/login/captcha", middleware.LoginRateLimitMiddleware(), systemHandler.GetCaptcha)
		systemRouter.POST("/user/login", middleware.LoginRateLimitMiddleware(), userHandler.Login)

		// 用户管理路由
		systemRouter.GET("/user/info", middleware.TokenVerify, userHandler.GetUserInfo)
		systemRouter.PUT("/user/info", middleware.TokenVerify, userHandler.UpdateUserInfo)
		systemRouter.GET("/user", middleware.TokenVerify, userHandler.FindUser)
		systemRouter.POST("/user", middleware.TokenVerify, userHandler.AddUser)
		systemRouter.PUT("/user", middleware.TokenVerify, userHandler.UpdateUser)
		systemRouter.DELETE("/user", middleware.TokenVerify, userHandler.DeleteUser)

		// 租户管理路由（需要超级管理员权限）
		systemRouter.GET("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenantHandler.FindTenant)
		systemRouter.POST("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenantHandler.AddTenant)
		systemRouter.PUT("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenantHandler.UpdateTenant)
		systemRouter.DELETE("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenantHandler.DeleteTenant)

		// TODO: 其他模块路由（菜单、部门、角色）需要按照相同模式重构
		// 暂时保留原有路由结构，等待后续完善
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
