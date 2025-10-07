package role

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func GetRoleList(c *gin.Context) {
	params := &struct {
		Name   string `json:"name" form:"name"`
		Status uint   `json:"status" form:"status"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取分页参数
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	isSuperAdmin := middleware.IsSuperAdmin(c)
	currentTenantID := middleware.GetTenantID(c)

	var targetTenantID uint
	if isSuperAdmin {
		tenantIDParam := c.Query("tenant_id")
		if tenantIDParam == "" {
			response.ReturnError(c, response.INVALID_ARGUMENT, "tenant_id 为必填参数")
			return
		}
		idValue, err := strconv.ParseUint(tenantIDParam, 10, 64)
		if err != nil || idValue == 0 {
			response.ReturnError(c, response.INVALID_ARGUMENT, "tenant_id 参数无效")
			return
		}
		targetTenantID = uint(idValue)
	} else {
		if currentTenantID == 0 {
			response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
			return
		}
		targetTenantID = currentTenantID
	}

	role := system.SystemRole{
		TenantID: targetTenantID,
		Name:     params.Name,
		Status:   params.Status,
	}

	// 调用带分页的查询函数
	roles, total, err := system.FindRoleList(&role, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取角色列表失败")
		return
	}

	if !isSuperAdmin {
		roleScopeIDs, err := system.GetTenantRoleScopeIDs(targetTenantID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取角色范围失败")
			return
		}
		if len(roleScopeIDs) > 0 {
			scopeSet := make(map[uint]struct{}, len(roleScopeIDs))
			for _, id := range roleScopeIDs {
				scopeSet[id] = struct{}{}
			}
			filtered := make([]system.SystemRole, 0, len(roles))
			for _, r := range roles {
				if _, ok := scopeSet[r.ID]; ok {
					filtered = append(filtered, r)
				}
			}
			roles = filtered
			total = int64(len(filtered))
		}
	}

	// 返回带总数的结果
	response.ReturnDataWithTotal(c, int(total), roles)
}

func AddRole(c *gin.Context) {
	params := &struct {
		TenantID uint   `json:"tenant_id"`
		Name     string `json:"name" form:"name" binding:"required"`
		Status   int    `json:"status" form:"status" binding:"required"`
		Desc     string `json:"desc" form:"desc"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if !middleware.IsSuperAdmin(c) {
		response.ReturnError(c, response.PERMISSION_DENIED, "仅平台管理员可以创建角色")
		return
	}
	if params.TenantID == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "tenant_id 为必填参数")
		return
	}
	role := system.SystemRole{
		TenantID: params.TenantID,
		Name:     params.Name,
		Status:   uint(params.Status),
		Desc:     params.Desc,
	}
	err := system.AddRole(&role)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加角色失败")
		return
	}
	if err := system.AddTenantRoleScope(params.TenantID, role.ID); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "同步角色范围失败")
		return
	}
	response.ReturnData(c, role)
}

func UpdateRole(c *gin.Context) {
	params := &struct {
		ID       uint   `json:"id" form:"id" binding:"required"`
		TenantID uint   `json:"tenant_id"`
		Name     string `json:"name" form:"name" binding:"required"`
		Status   int    `json:"status" form:"status" binding:"required"`
		Desc     string `json:"desc" form:"desc"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	isSuperAdmin := middleware.IsSuperAdmin(c)
	currentTenantID := middleware.GetTenantID(c)

	originalRole := system.SystemRole{Model: gorm.Model{ID: params.ID}}
	if err := system.GetRole(&originalRole); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "角色不存在")
		return
	}

	targetTenantID := originalRole.TenantID
	if isSuperAdmin {
		if params.TenantID != 0 && params.TenantID != originalRole.TenantID {
			targetTenantID = params.TenantID
		}
	} else {
		if currentTenantID == 0 || originalRole.TenantID != currentTenantID {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权操作该角色")
			return
		}
		roleScopeIDs, err := system.GetTenantRoleScopeIDs(currentTenantID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取角色范围失败")
			return
		}
		if len(roleScopeIDs) > 0 {
			allowed := false
			for _, id := range roleScopeIDs {
				if id == originalRole.ID {
					allowed = true
					break
				}
			}
			if !allowed {
				response.ReturnError(c, response.PERMISSION_DENIED, "角色不在可管理范围内")
				return
			}
		}
	}

	updatedRole := system.SystemRole{
		Model:    gorm.Model{ID: params.ID},
		TenantID: targetTenantID,
		Name:     params.Name,
		Status:   uint(params.Status),
		Desc:     params.Desc,
	}
	err := system.UpdateRole(&updatedRole)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新角色失败")
		return
	}
	if isSuperAdmin && targetTenantID != originalRole.TenantID {
		// 移动角色到新的租户后，需要更新范围
		if err := system.RemoveTenantRoleScope(originalRole.TenantID, updatedRole.ID); err != nil {
			response.ReturnError(c, response.DATA_LOSS, "同步角色范围失败")
			return
		}
		if err := system.AddTenantRoleScope(targetTenantID, updatedRole.ID); err != nil {
			response.ReturnError(c, response.DATA_LOSS, "同步角色范围失败")
			return
		}
	}
	response.ReturnData(c, updatedRole)
}

func DeleteRole(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	role := system.SystemRole{Model: gorm.Model{ID: params.ID}}
	if err := system.GetRole(&role); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "角色不存在")
		return
	}
	if !middleware.IsSuperAdmin(c) {
		tenantID := middleware.GetTenantID(c)
		if tenantID == 0 || role.TenantID != tenantID {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权删除该角色")
			return
		}
		roleScopeIDs, err := system.GetTenantRoleScopeIDs(tenantID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取角色范围失败")
			return
		}
		if len(roleScopeIDs) > 0 {
			allowed := false
			for _, id := range roleScopeIDs {
				if id == role.ID {
					allowed = true
					break
				}
			}
			if !allowed {
				response.ReturnError(c, response.PERMISSION_DENIED, "角色不在可管理范围内")
				return
			}
		}
	}
	if err := system.DeleteRole(&role); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除角色失败")
		return
	}
	if middleware.IsSuperAdmin(c) {
		if err := system.RemoveTenantRoleScope(role.TenantID, role.ID); err != nil {
			response.ReturnError(c, response.DATA_LOSS, "同步角色范围失败")
			return
		}
	}
	response.ReturnData(c, role)
}
