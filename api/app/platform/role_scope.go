package platform

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func GetTenantRoleScope(c *gin.Context) {
	tenantIDParam := c.Query("tenant_id")
	if tenantIDParam == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "tenant_id 为必填参数")
		return
	}
	tenantIDValue, err := strconv.ParseUint(tenantIDParam, 10, 64)
	if err != nil || tenantIDValue == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "tenant_id 参数无效")
		return
	}
	roleIDs, err := system.GetTenantRoleScopeIDs(uint(tenantIDValue))
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取角色范围失败")
		return
	}
	response.ReturnData(c, gin.H{
		"tenant_id": uint(tenantIDValue),
		"role_ids":  roleIDs,
	})
}

func UpdateTenantRoleScope(c *gin.Context) {
	params := &struct {
		TenantID uint   `json:"tenant_id" binding:"required"`
		RoleIDs  []uint `json:"role_ids"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if err := system.SaveTenantRoleScope(params.TenantID, params.RoleIDs); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "保存角色范围失败")
		return
	}
	response.ReturnData(c, gin.H{
		"tenant_id": params.TenantID,
		"role_ids":  params.RoleIDs,
	})
}
