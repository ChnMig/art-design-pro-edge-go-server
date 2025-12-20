package menu

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	commonmenu "api-server/common/menu"
	menudomain "api-server/domain/admin/menu"
)

// GetMenuListByRoleID 根据角色ID获取菜单列表
func GetMenuListByRoleID(c *gin.Context) {
	params := &struct {
		RoleID uint `json:"role_id" form:"role_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	menuTree, err := menudomain.GetRoleMenuTree(params.RoleID, middleware.GetTenantID(c), middleware.IsSuperAdmin(c))
	if err != nil {
		if errors.Is(err, menudomain.ErrPermissionDenied) {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权查看该角色菜单")
			return
		}
		ReturnDomainError(c, err, "查询角色菜单失败")
		return
	}
	response.ReturnData(c, menuTree)
}

func UpdateMenuListByRoleID(c *gin.Context) {
	params := &struct {
		RoleID   uint   `json:"role_id" form:"role_id" binding:"required"`
		MenuData string `json:"menu_data" form:"menu_data" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 尝试将 params.MenuData 转成结构体
	var menuData []commonmenu.MenuResponse
	err := json.Unmarshal([]byte(params.MenuData), &menuData)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "参数错误")
		return
	}

	if err := menudomain.UpdateRoleMenu(params.RoleID, menuData, middleware.GetTenantID(c), middleware.IsSuperAdmin(c)); err != nil {
		if errors.Is(err, menudomain.ErrPermissionDenied) {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权调整该角色菜单")
			return
		}
		ReturnDomainError(c, err, "保存角色菜单失败")
		return
	}

	response.ReturnData(c, nil)
}
