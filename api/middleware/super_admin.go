package middleware

import (
	"github.com/gin-gonic/gin"

	"api-server/api/response"
	"api-server/db/pgdb/system"
)

// SuperAdminVerify 超级管理员权限验证中间件
// 只有用户ID为1的超级管理员才能执行租户管理操作
func SuperAdminVerify(c *gin.Context) {
	// 获取当前用户ID
	userID := GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户认证失败")
		c.Abort()
		return
	}

	// 检查是否为超级管理员（用户ID为1）
	if userID != 1 {
		response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，只有超级管理员可以管理租户")
		c.Abort()
		return
	}

	// 验证用户存在且状态正常
	user := system.SystemUser{}
	user.ID = userID
	if err := system.GetUser(&user); err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户不存在")
		c.Abort()
		return
	}

	if user.Status != system.StatusEnabled {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户已被禁用")
		c.Abort()
		return
	}

	c.Next()
}

// TenantAdminVerify 租户管理员权限验证中间件
// 允许超级管理员或租户管理员执行特定操作
func TenantAdminVerify(c *gin.Context) {
	// 获取当前用户ID和租户ID
	userID := GetCurrentUserID(c)
	tenantID := GetTenantID(c)

	if userID == 0 || tenantID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户认证失败")
		c.Abort()
		return
	}

	// 超级管理员（用户ID为1）可以访问任何租户数据
	if userID == 1 {
		c.Next()
		return
	}

	// 验证用户存在且状态正常
	user := system.SystemUser{}
	user.ID = userID
	user.TenantID = tenantID
	if err := system.GetUser(&user); err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户不存在或不属于指定租户")
		c.Abort()
		return
	}

	if user.Status != system.StatusEnabled {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户已被禁用")
		c.Abort()
		return
	}

	// 检查用户角色ID（假设角色ID为1是租户管理员）
	if user.RoleID != 1 && user.RoleID != 2 { // 1: 超级管理员, 2: 租户管理员
		response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，需要管理员权限")
		c.Abort()
		return
	}

	c.Next()
}
