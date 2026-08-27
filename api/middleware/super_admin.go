package middleware

import (
	"github.com/gin-gonic/gin"

	"api-server/api/response"
	"api-server/db/pgdb/system"
)

// SuperAdminVerify 超级管理员权限验证中间件
// 只有用户ID为1的超级管理员才能执行租户管理操作
func SuperAdminVerify(c *gin.Context) {
	user, ok := GetCurrentUser(c)
	tenant, tenantOK := GetCurrentTenant(c)
	tenantID := GetTenantID(c)
	if !ok || !tenantOK || user.ID == 0 || tenantID == 0 ||
		user.TenantID != tenantID || tenant.ID != tenantID ||
		user.Status != system.StatusEnabled || tenant.Status != system.StatusEnabled {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户认证失败")
		c.Abort()
		return
	}

	// 检查是否为超级管理员（用户ID为1）
	if user.ID != 1 {
		response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，只有超级管理员可以管理租户")
		c.Abort()
		return
	}

	c.Next()
}

// IsSuperAdmin 判断当前请求是否为平台超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	user, ok := GetCurrentUser(c)
	return ok && user.ID == 1 && user.Status == system.StatusEnabled
}

// TenantAdminVerify 租户管理员权限验证中间件
// 允许超级管理员或租户管理员执行特定操作
func TenantAdminVerify(c *gin.Context) {
	user, ok := GetCurrentUser(c)
	tenant, tenantOK := GetCurrentTenant(c)
	tenantID := GetTenantID(c)

	if !ok || !tenantOK || user.ID == 0 || tenantID == 0 ||
		user.TenantID != tenantID || tenant.ID != tenantID ||
		user.Status != system.StatusEnabled || tenant.Status != system.StatusEnabled {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户认证失败")
		c.Abort()
		return
	}

	// 超级管理员（用户ID为1）可以访问任何租户数据
	if user.ID == 1 {
		c.Next()
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
