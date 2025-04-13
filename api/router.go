package api

import (
	"github.com/gin-gonic/gin"

	"api-server/api/app/system/department"
	"api-server/api/app/system/menu"
	"api-server/api/app/system/role"
	"api-server/api/app/system/todo"
	"api-server/api/app/system/user"
	"api-server/api/middleware"
)

func systemRouter(router *gin.RouterGroup) {
	systemRouter := router.Group("/admin/system")
	{
		systemRouter.GET("/user/login/captcha", user.GetCaptcha)
		systemRouter.POST("/user/login", user.Login)
		systemRouter.GET("/user/info", middleware.TokenVerify, user.GetUserInfo)
		systemRouter.PUT("/user/info", middleware.TokenVerify, user.UpdateUserInfo)
		systemRouter.GET("/user/menu", middleware.TokenVerify, user.GetUserMenuList)
		systemRouter.GET("/menu", middleware.TokenVerify, menu.GetMenuList)
		systemRouter.POST("/menu", middleware.TokenVerify, menu.AddMenu)
		systemRouter.DELETE("/menu", middleware.TokenVerify, menu.DeleteMenu)
		systemRouter.PUT("/menu", middleware.TokenVerify, menu.UpdateMenu)
		systemRouter.GET("/menu/auth", middleware.TokenVerify, menu.GetMenuAuthList)
		systemRouter.POST("/menu/auth", middleware.TokenVerify, menu.AddMenuAuth)
		systemRouter.DELETE("/menu/auth", middleware.TokenVerify, menu.DeleteMenuAuth)
		systemRouter.PUT("/menu/auth", middleware.TokenVerify, menu.UpdateMenuAuth)
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

		// Todo相关路由
		systemRouter.GET("/todo", middleware.TokenVerify, todo.FindTodoList)
		systemRouter.GET("/todo/detail", middleware.TokenVerify, todo.GetTodo)
		systemRouter.POST("/todo", middleware.TokenVerify, todo.AddTodo)
		systemRouter.PUT("/todo/status", middleware.TokenVerify, todo.UpdateTodoStatus)
		systemRouter.PUT("/todo", middleware.TokenVerify, todo.UpdateTodo) // 新增路由，用于更新待办事项的全部信息
		systemRouter.DELETE("/todo", middleware.TokenVerify, todo.DeleteTodo)

		// Todo评论相关路由
		systemRouter.GET("/todo/comment", middleware.TokenVerify, todo.FindTodoComments)
		systemRouter.POST("/todo/comment", middleware.TokenVerify, todo.AddTodoComment)

		// Todo步骤相关路由
		systemRouter.GET("/todo/step", middleware.TokenVerify, todo.FindTodoSteps)
		systemRouter.POST("/todo/step", middleware.TokenVerify, todo.AddTodoStep)

		// Todo日志相关路由
		systemRouter.GET("/todo/log", middleware.TokenVerify, todo.FindTodoLogs)
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
