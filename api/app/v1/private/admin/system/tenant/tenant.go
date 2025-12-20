package tenant

import (
	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	tenantdomain "api-server/domain/admin/tenant"
)

// FindTenant 查询租户列表
func FindTenant(c *gin.Context) {
	params := &struct {
		Code   string `json:"code" form:"code"`
		Name   string `json:"name" form:"name"`
		Status uint   `json:"status" form:"status"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	tenants, total, err := tenantdomain.FindTenantList(tenantdomain.FindListQuery{
		Code:   params.Code,
		Name:   params.Name,
		Status: params.Status,
	}, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询租户失败")
		return
	}

	response.ReturnDataWithTotal(c, int(total), tenants)
}

// AddTenant 添加租户
func AddTenant(c *gin.Context) {
	params := &struct {
		Code    string `json:"code" form:"code" binding:"required"`
		Name    string `json:"name" form:"name" binding:"required"`
		Contact string `json:"contact" form:"contact"`
		Phone   string `json:"phone" form:"phone"`
		Email   string `json:"email" form:"email"`
		Status  uint   `json:"status" form:"status" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	_, err := tenantdomain.AddTenant(tenantdomain.AddTenantInput{
		Code:    params.Code,
		Name:    params.Name,
		Contact: params.Contact,
		Phone:   params.Phone,
		Email:   params.Email,
		Status:  params.Status,
	})
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加租户失败")
		return
	}

	response.ReturnData(c, nil)
}

// UpdateTenant 更新租户
func UpdateTenant(c *gin.Context) {
	params := &struct {
		ID      uint   `json:"id" form:"id" binding:"required"`
		Code    string `json:"code" form:"code" binding:"required"`
		Name    string `json:"name" form:"name" binding:"required"`
		Contact string `json:"contact" form:"contact"`
		Phone   string `json:"phone" form:"phone"`
		Email   string `json:"email" form:"email"`
		Status  uint   `json:"status" form:"status" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	if err := tenantdomain.UpdateTenant(tenantdomain.UpdateTenantInput{
		ID:      params.ID,
		Code:    params.Code,
		Name:    params.Name,
		Contact: params.Contact,
		Phone:   params.Phone,
		Email:   params.Email,
		Status:  params.Status,
	}); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新租户失败")
		return
	}

	response.ReturnData(c, nil)
}

// DeleteTenant 删除租户
func DeleteTenant(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	if err := tenantdomain.DeleteTenant(params.ID); err != nil {
		ReturnDomainError(c, err, "删除租户失败")
		return
	}

	response.ReturnData(c, nil)
}
