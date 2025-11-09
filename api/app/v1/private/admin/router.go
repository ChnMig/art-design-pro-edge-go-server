package admin

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
)

// RegisterRoutes 在 /api/v1/private/admin 下注册系统管理接口
func RegisterRoutes(admin *gin.RouterGroup) {
	if admin == nil {
		return
	}

	registerSystemRoutes(admin.Group("/system"))
	registerPlatformRoutes(admin.Group("/platform"))
}

func registerSystemRoutes(group *gin.RouterGroup) {
	if group == nil {
		return
	}

	group.GET("/user/login/captcha", middleware.LoginRateLimitMiddleware(), user.GetCaptcha)
	group.POST("/user/login", middleware.LoginRateLimitMiddleware(), user.Login)
	group.GET("/user/login/tenant", middleware.LoginRateLimitMiddleware(), user.SearchTenantCodeForLogin)
	group.GET("/login/log", middleware.TokenVerify, user.FindLoginLogList)
	group.GET("/user/info", middleware.TokenVerify, user.GetUserInfo)
	group.PUT("/user/info", middleware.TokenVerify, user.UpdateUserInfo)
	group.GET("/user/menu", middleware.TokenVerify, user.GetUserMenuList)
	group.GET("/menu/role", middleware.TokenVerify, menu.GetMenuListByRoleID)
	group.PUT("/menu/role", middleware.TokenVerify, menu.UpdateMenuListByRoleID)
	group.GET("/department", middleware.TokenVerify, department.GetDepartmentList)
	group.POST("/department", middleware.TokenVerify, department.AddDepartment)
	group.PUT("/department", middleware.TokenVerify, department.UpdateDepartment)
	group.DELETE("/department", middleware.TokenVerify, department.DeleteDepartment)
	group.GET("/role", middleware.TokenVerify, role.GetRoleList)
	group.POST("/role", middleware.TokenVerify, role.AddRole)
	group.PUT("/role", middleware.TokenVerify, role.UpdateRole)
	group.DELETE("/role", middleware.TokenVerify, role.DeleteRole)
	group.GET("/user", middleware.TokenVerify, user.FindUser)
	group.GET("/user/cache", middleware.TokenVerify, user.FindUserByCache)
	group.POST("/user", middleware.TokenVerify, user.AddUser)
	group.PUT("/user", middleware.TokenVerify, user.UpdateUser)
	group.DELETE("/user", middleware.TokenVerify, user.DeleteUser)
	group.GET("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.FindTenant)
	group.POST("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.AddTenant)
	group.PUT("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.UpdateTenant)
	group.DELETE("/tenant", middleware.TokenVerify, middleware.SuperAdminVerify, tenant.DeleteTenant)
}

func registerPlatformRoutes(group *gin.RouterGroup) {
	if group == nil {
		return
	}

	group.Use(middleware.TokenVerify, middleware.SuperAdminVerify)
	group.GET("/menu", platformMenu.GetMenuList)
	group.POST("/menu", platformMenu.AddMenu)
	group.PUT("/menu", platformMenu.UpdateMenu)
	group.DELETE("/menu", platformMenu.DeleteMenu)
	group.GET("/menu/tenant", platformMenu.GetTenantMenu)
	group.PUT("/menu/tenant", platformMenu.UpdateTenantMenu)
	group.GET("/menu/auth", platformMenu.GetMenuAuthList)
	group.POST("/menu/auth", platformMenu.AddMenuAuth)
	group.PUT("/menu/auth", platformMenu.UpdateMenuAuth)
	group.DELETE("/menu/auth", platformMenu.DeleteMenuAuth)

	group.GET("/role", platformRole.GetRoleList)
	group.POST("/role", platformRole.AddRole)
	group.PUT("/role", platformRole.UpdateRole)
	group.DELETE("/role", platformRole.DeleteRole)
	group.GET("/tenant", tenant.FindTenant)
	group.POST("/tenant", tenant.AddTenant)
	group.PUT("/tenant", tenant.UpdateTenant)
	group.DELETE("/tenant", tenant.DeleteTenant)
}
