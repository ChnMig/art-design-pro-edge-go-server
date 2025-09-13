package system

import (
	"github.com/gin-gonic/gin"

	"api-server/internal/service/system"
	"api-server/internal/transport/http/middleware"
	"api-server/internal/transport/http/response"
)

// TenantHandler 租户HTTP处理器
type TenantHandler struct {
	tenantService system.TenantService
}

// NewTenantHandler 创建租户处理器
func NewTenantHandler(tenantService system.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// FindTenant 查询租户列表处理器
func (h *TenantHandler) FindTenant(c *gin.Context) {
	// 1. 参数绑定
	req := &system.FindTenantListRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 获取分页参数
	req.Page = middleware.GetPage(c)
	req.PageSize = middleware.GetPageSize(c)

	// 3. 调用租户服务
	result, err := h.tenantService.FindTenantList(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 4. 返回结果
	response.ReturnDataWithTotal(c, int(result.Total), result.Tenants)
}

// AddTenant 添加租户处理器
func (h *TenantHandler) AddTenant(c *gin.Context) {
	// 1. 参数绑定和验证
	req := &system.CreateTenantRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 调用租户服务
	err := h.tenantService.CreateTenant(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 3. 返回成功
	response.ReturnData(c, nil)
}

// UpdateTenant 更新租户处理器
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	// 1. 参数绑定和验证
	req := &system.UpdateTenantRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 调用租户服务
	err := h.tenantService.UpdateTenant(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 3. 返回成功
	response.ReturnData(c, nil)
}

// DeleteTenant 删除租户处理器
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	// 1. 参数绑定
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 2. 调用租户服务
	err := h.tenantService.DeleteTenant(c.Request.Context(), params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 3. 返回成功
	response.ReturnData(c, nil)
}