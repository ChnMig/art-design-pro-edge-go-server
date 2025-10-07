package tenant

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
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
	tenant := system.SystemTenant{
		Code:   params.Code,
		Name:   params.Name,
		Status: params.Status,
	}

	tenants, total, err := system.FindTenantList(&tenant, page, pageSize)
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
		Address string `json:"address" form:"address"`
		Status  uint   `json:"status" form:"status" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	tenant := system.SystemTenant{
		Code:    params.Code,
		Name:    params.Name,
		Contact: params.Contact,
		Phone:   params.Phone,
		Email:   params.Email,
		Address: params.Address,
		Status:  params.Status,
	}

	if err := system.AddTenant(&tenant); err != nil {
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
		Address string `json:"address" form:"address"`
		Status  uint   `json:"status" form:"status" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	tenant := system.SystemTenant{
		Model:   gorm.Model{ID: params.ID},
		Code:    params.Code,
		Name:    params.Name,
		Contact: params.Contact,
		Phone:   params.Phone,
		Email:   params.Email,
		Address: params.Address,
		Status:  params.Status,
	}

	if err := system.UpdateTenant(&tenant); err != nil {
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

	tenant := system.SystemTenant{
		Model: gorm.Model{ID: params.ID},
	}

	if err := system.DeleteTenant(&tenant); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除租户失败")
		return
	}

	response.ReturnData(c, nil)
}
