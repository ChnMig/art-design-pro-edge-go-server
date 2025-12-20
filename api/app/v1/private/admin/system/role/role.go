package role

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	roledomain "api-server/domain/admin/role"
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

	// 从 JWT 中获取租户 ID，系统侧不再需要前端传 tenant_id
	currentTenantID := middleware.GetTenantID(c)
	if currentTenantID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}
	targetTenantID := currentTenantID

	// 调用带分页的查询函数
	roles, total, err := roledomain.FindRoleList(roledomain.FindListQuery{
		TenantID: targetTenantID,
		Name:     params.Name,
		Status:   params.Status,
	}, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取角色列表失败")
		return
	}

	// 返回带总数的结果
	response.ReturnDataWithTotal(c, int(total), roles)
}

func AddRole(c *gin.Context) {
	params := &struct {
		Name   string `json:"name" form:"name" binding:"required"`
		Status int    `json:"status" form:"status" binding:"required"`
		Desc   string `json:"desc" form:"desc"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	tenantID := middleware.GetTenantID(c)
	if !middleware.IsSuperAdmin(c) {
		if tenantID == 0 {
			response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
			return
		}
	}

	targetID := tenantID
	if middleware.IsSuperAdmin(c) {
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
		targetID = uint(idValue)
	}

	role, err := roledomain.AddRole(roledomain.AddInput{
		TenantID: targetID,
		Name:     params.Name,
		Status:   uint(params.Status),
		Desc:     params.Desc,
	})
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加角色失败")
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

	originalRole, err := roledomain.GetRole(params.ID)
	if err != nil {
		ReturnDomainError(c, err, "角色不存在")
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
		// 非超级管理员不允许调整角色归属租户
		targetTenantID = originalRole.TenantID
	}

	updatedRole, err := roledomain.UpdateRole(roledomain.UpdateInput{
		ID:       params.ID,
		TenantID: targetTenantID,
		Name:     params.Name,
		Status:   uint(params.Status),
		Desc:     params.Desc,
	})
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新角色失败")
		return
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
	roleEntity, err := roledomain.GetRole(params.ID)
	if err != nil {
		ReturnDomainError(c, err, "角色不存在")
		return
	}
	if !middleware.IsSuperAdmin(c) {
		tenantID := middleware.GetTenantID(c)
		if tenantID == 0 || roleEntity.TenantID != tenantID {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权删除该角色")
			return
		}
	}
	if err := roledomain.DeleteRole(params.ID); err != nil {
		ReturnDomainError(c, err, "删除角色失败")
		return
	}
	response.ReturnData(c, roleEntity)
}
