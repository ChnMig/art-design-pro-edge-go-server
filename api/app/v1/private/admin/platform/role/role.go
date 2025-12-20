package role

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	roledomain "api-server/domain/admin/role"
)

func GetRoleList(c *gin.Context) {
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
	params := &struct {
		Name   string `form:"name"`
		Status uint   `form:"status"`
	}{}
	if err := c.ShouldBindQuery(params); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "查询参数无效")
		return
	}

	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	roles, total, err := roledomain.FindRoleList(roledomain.FindListQuery{
		TenantID: uint(tenantIDValue),
		Name:     params.Name,
		Status:   params.Status,
	}, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取角色列表失败")
		return
	}
	response.ReturnDataWithTotal(c, int(total), roles)
}

func AddRole(c *gin.Context) {
	params := &struct {
		TenantID uint   `json:"tenant_id" binding:"required"`
		Name     string `json:"name" form:"name" binding:"required"`
		Status   int    `json:"status" form:"status" binding:"required"`
		Desc     string `json:"desc" form:"desc"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	role, err := roledomain.AddRole(roledomain.AddInput{
		TenantID: params.TenantID,
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
	role, err := roledomain.UpdateRole(roledomain.UpdateInput{
		ID:       params.ID,
		TenantID: params.TenantID,
		Name:     params.Name,
		Status:   uint(params.Status),
		Desc:     params.Desc,
	})
	if err != nil {
		ReturnDomainError(c, err, "更新角色失败")
		return
	}
	response.ReturnData(c, role)
}

func DeleteRole(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	role, err := roledomain.GetRole(params.ID)
	if err != nil {
		ReturnDomainError(c, err, "角色不存在")
		return
	}
	if err := roledomain.DeleteRole(params.ID); err != nil {
		ReturnDomainError(c, err, "删除角色失败")
		return
	}
	response.ReturnData(c, role)
}
