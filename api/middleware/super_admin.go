package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/response"
	"api-server/db/pgdb/system"
)

type roleLoader func(role *system.SystemRole) error

// SuperAdminVerify 平台超级管理员权限验证中间件。
// 只有默认平台租户中用户 ID 为 1 的启用用户才能执行平台管理操作。
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

	// 复核平台根用户的固定身份边界。
	if !isPlatformSuperAdmin(user, tenantID) {
		response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，只有超级管理员可以管理租户")
		c.Abort()
		return
	}

	c.Next()
}

// IsSuperAdmin 判断当前请求是否为平台超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	user, ok := GetCurrentUser(c)
	return ok && isPlatformSuperAdmin(user, GetTenantID(c))
}

// TenantAdminVerify 租户管理员权限验证中间件
// 允许超级管理员或租户管理员执行特定操作
func TenantAdminVerify(c *gin.Context) {
	tenantAdminVerify(c, system.GetRole)
}

func tenantAdminVerify(c *gin.Context, loadRole roleLoader) {
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

	// 平台根用户可进入平台级管理流程，具体目标租户仍由接口参数显式指定。
	if isPlatformSuperAdmin(user, tenantID) {
		c.Next()
		return
	}

	if user.RoleID == 0 || loadRole == nil {
		response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，需要管理员权限")
		c.Abort()
		return
	}
	role := system.SystemRole{Model: gorm.Model{ID: user.RoleID}, TenantID: tenantID}
	if err := loadRole(&role); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.UNAUTHENTICATED, "角色状态校验失败")
		} else {
			response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，需要管理员权限")
		}
		c.Abort()
		return
	}
	if role.TenantID != tenantID || role.Status != system.StatusEnabled || !role.IsTenantAdmin {
		response.ReturnError(c, response.PERMISSION_DENIED, "权限不足，需要管理员权限")
		c.Abort()
		return
	}

	c.Next()
}

func isPlatformSuperAdmin(user system.SystemUser, tenantID uint) bool {
	return user.ID == 1 && user.TenantID == system.DefaultTenantID && tenantID == system.DefaultTenantID && user.Status == system.StatusEnabled
}
