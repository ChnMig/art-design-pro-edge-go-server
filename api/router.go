package api

import (
	"github.com/gin-gonic/gin"

	platformMenu "api-server/api/app/platform/menu"
	platformRole "api-server/api/app/platform/role"
	"api-server/api/app/system/department"
	"api-server/api/app/system/menu"
	"api-server/api/app/system/role"
	"api-server/api/app/system/tenant"
	"api-server/api/app/system/user"
	"api-server/api/middleware"
	"api-server/config"
)

func systemRouter(router *gin.RouterGroup) {
	systemRouter := router.Group("/admin/system")
	{
		systemRouter.GET("/user/login/captcha", middleware.LoginRateLimitMiddleware(), user.GetCaptcha)
		systemRouter.POST("/user/login", middleware.LoginRateLimitMiddleware(), user.Login)
		systemRouter.GET("/user/login/tenant", middleware.LoginRateLimitMiddleware(), user.SearchTenantCodeForLogin)
		systemRouter.GET("/login/log", middleware.TokenVerify, user.FindLoginLogList)
		systemRouter.GET("/user/info", middleware.TokenVerify, user.GetUserInfo)
		systemRouter.PUT("/user/info", middleware.TokenVerify, user.UpdateUserInfo)
		systemRouter.GET("/user/menu", middleware.TokenVerify, user.GetUserMenuList)
		systemRouter.GET("/menu/role", middleware.TokenVerify, menu.GetMenuListByRoleID)
		systemRouter.PUT("/menu/role", middleware.TokenVerify, menu.UpdateMenuListByRoleID)
		systemRouter.GET("/department", middleware.TokenVerify, department.GetDepartmentList)
		systemRouter.POST("/department", middleware.TokenVerify, department.AddDepartment)
		systemRouter.PUT("/department", middleware.TokenVerify, department.UpdateDepartment)
		systemRouter.DELETE("/department", middleware.TokenVerify, department.DeleteDepartment)
		systemRouter.GET("/role", middleware.TokenVerify, role.GetRoleList)
		systemRouter.POST("/role", middleware.TokenVerify, role.AddRole)
		systemRouter.PUT("/role", middleware.TokenVerify, role.UpdateRole)
		systemRouter.DELETE("/role", middleware.TokenVerify, role.DeleteRole)
		systemRouter.GET("/user", middleware.TokenVerify, user.FindUser)
		systemRouter.GET("/user/cache", middleware.TokenVerify, user.FindUserByCache)
		systemRouter.POST("/user", middleware.TokenVerify, user.AddUser)
		systemRouter.PUT("/user", middleware.TokenVerify, user.UpdateUser)
		systemRouter.DELETE("/user", middleware.TokenVerify, user.DeleteUser)
		systemRouter.GET("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.FindTenant)
		systemRouter.POST("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.AddTenant)
		systemRouter.PUT("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.UpdateTenant)
		systemRouter.DELETE("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.DeleteTenant)
	}

	platformRouter := router.Group("/admin/platform", middleware.TokenVerify, middleware.SuperAdminVerify)
	{
		platformRouter.GET("/menu", platformMenu.GetMenuList)
		platformRouter.POST("/menu", platformMenu.AddMenu)
		platformRouter.PUT("/menu", platformMenu.UpdateMenu)
		platformRouter.DELETE("/menu", platformMenu.DeleteMenu)
		// 租户菜单范围
		platformRouter.GET("/menu/tenant", platformMenu.GetTenantMenu)
		platformRouter.PUT("/menu/tenant", platformMenu.UpdateTenantMenu)
		platformRouter.GET("/menu/auth", platformMenu.GetMenuAuthList)
		platformRouter.POST("/menu/auth", platformMenu.AddMenuAuth)
		platformRouter.PUT("/menu/auth", platformMenu.UpdateMenuAuth)
		platformRouter.DELETE("/menu/auth", platformMenu.DeleteMenuAuth)

		platformRouter.GET("/role", platformRole.GetRoleList)
		platformRouter.POST("/role", platformRole.AddRole)
		platformRouter.PUT("/role", platformRole.UpdateRole)
		platformRouter.DELETE("/role", platformRole.DeleteRole)
		platformRouter.GET("/tenant", tenant.FindTenant)
		platformRouter.POST("/tenant", tenant.AddTenant)
		platformRouter.PUT("/tenant", tenant.UpdateTenant)
		platformRouter.DELETE("/tenant", tenant.DeleteTenant)
	}
}

// InitApi init gshop app
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

	// static
	router.Static("/static", "./static")
	// api-v1
	v1 := router.Group("/api/v1")
	{
		systemRouter(v1)
	}
	return router
}
