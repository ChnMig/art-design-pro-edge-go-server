package menu

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	commonmenu "api-server/common/menu"
	menudomain "api-server/domain/admin/menu"
)

// GetMenuList 获取平台菜单定义（不带租户 hasPermission 标记）。
// GET /api/v1/admin/platform/menu
func GetMenuList(c *gin.Context) {
	menuTree, err := menudomain.GetPlatformMenuTree()
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	response.ReturnData(c, menuTree)
}

func AddMenu(c *gin.Context) {
	params := &struct {
		Path          string `json:"path" form:"path" binding:"required"`
		Name          string `json:"name" form:"name" binding:"required"`
		Component     string `json:"component" form:"component"`
		Title         string `json:"title" form:"title" binding:"required"`
		Icon          string `json:"icon" form:"icon"`
		ShowBadge     uint   `json:"showBadge" form:"showBadge"`
		ShowTextBadge string `json:"showTextBadge" form:"showTextBadge"`
		IsHide        uint   `json:"isHide" form:"isHide" binding:"required"`
		IsHideTab     uint   `json:"isHideTab" form:"isHideTab" binding:"required"`
		Link          string `json:"link" form:"link"`
		IsIframe      uint   `json:"isIframe" form:"isIframe" binding:"required"`
		KeepAlive     uint   `json:"keepAlive" form:"keepAlive" binding:"required"`
		IsFirstLevel  uint   `json:"isFirstLevel" form:"isFirstLevel" binding:"required"`
		Status        uint   `json:"status" form:"status" binding:"required"`
		ParentID      uint   `json:"parentId" form:"parentId"`
		Sort          uint   `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	menuEntity, err := menudomain.AddMenu(menudomain.AddMenuInput{
		Path:          params.Path,
		Name:          params.Name,
		Component:     params.Component,
		Title:         params.Title,
		Icon:          params.Icon,
		ShowBadge:     params.ShowBadge,
		ShowTextBadge: params.ShowTextBadge,
		IsHide:        params.IsHide,
		IsHideTab:     params.IsHideTab,
		Link:          params.Link,
		IsIframe:      params.IsIframe,
		KeepAlive:     params.KeepAlive,
		IsFirstLevel:  params.IsFirstLevel,
		Status:        params.Status,
		ParentID:      params.ParentID,
		Sort:          params.Sort,
	})
	if err != nil {
		ReturnDomainError(c, err, "添加菜单失败")
		return
	}
	response.ReturnData(c, menuEntity)
}

// UpdateMenu 更新平台菜单定义
// PUT /api/v1/admin/platform/menu
func UpdateMenu(c *gin.Context) {
	// 菜单定义更新
	params := &struct {
		ID            uint   `json:"id" form:"id" binding:"required"`
		Path          string `json:"path" form:"path" binding:"required"`
		Name          string `json:"name" form:"name" binding:"required"`
		Component     string `json:"component" form:"component"`
		Title         string `json:"title" form:"title" binding:"required"`
		Icon          string `json:"icon" form:"icon"`
		ShowBadge     uint   `json:"showBadge" form:"showBadge"`
		ShowTextBadge string `json:"showTextBadge" form:"showTextBadge"`
		IsHide        uint   `json:"isHide" form:"isHide" binding:"required"`
		IsHideTab     uint   `json:"isHideTab" form:"isHideTab" binding:"required"`
		Link          string `json:"link" form:"link"`
		IsIframe      uint   `json:"isIframe" form:"isIframe" binding:"required"`
		KeepAlive     uint   `json:"keepAlive" form:"keepAlive" binding:"required"`
		IsFirstLevel  uint   `json:"isFirstLevel" form:"isFirstLevel" binding:"required"`
		Status        uint   `json:"status" form:"status" binding:"required"`
		ParentID      uint   `json:"parentId" form:"parentId"`
		Sort          uint   `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	menuEntity, err := menudomain.UpdateMenu(menudomain.UpdateMenuInput{
		ID:            params.ID,
		Path:          params.Path,
		Name:          params.Name,
		Component:     params.Component,
		Title:         params.Title,
		Icon:          params.Icon,
		ShowBadge:     params.ShowBadge,
		ShowTextBadge: params.ShowTextBadge,
		IsHide:        params.IsHide,
		IsHideTab:     params.IsHideTab,
		Link:          params.Link,
		IsIframe:      params.IsIframe,
		KeepAlive:     params.KeepAlive,
		IsFirstLevel:  params.IsFirstLevel,
		Status:        params.Status,
		ParentID:      params.ParentID,
		Sort:          params.Sort,
	})
	if err != nil {
		ReturnDomainError(c, err, "更新菜单失败")
		return
	}
	response.ReturnData(c, menuEntity)
}

func DeleteMenu(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	menuEntity, err := menudomain.DeleteMenu(params.ID)
	if err != nil {
		ReturnDomainError(c, err, "删除菜单失败")
		return
	}
	response.ReturnData(c, menuEntity)
}

func GetMenuAuthList(c *gin.Context) {
	params := &struct {
		MenuID uint `json:"menu_id" form:"menu_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auths, err := menudomain.GetMenuAuthList(params.MenuID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单权限失败")
		return
	}
	response.ReturnData(c, auths)
}

func AddMenuAuth(c *gin.Context) {
	params := &struct {
		MenuID uint   `json:"menu_id"`
		Mark   string `json:"mark"`
		Title  string `json:"title"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	auth, err := menudomain.AddMenuAuth(menudomain.AddMenuAuthInput{
		MenuID: params.MenuID,
		Mark:   params.Mark,
		Title:  params.Title,
	})
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}

func UpdateMenuAuth(c *gin.Context) {
	params := &struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Mark   string `json:"mark"`
		MenuID uint   `json:"menu_id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	auth, err := menudomain.UpdateMenuAuth(menudomain.UpdateMenuAuthInput{
		ID:     params.ID,
		Title:  params.Title,
		Mark:   params.Mark,
		MenuID: params.MenuID,
	})
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}

func DeleteMenuAuth(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth, err := menudomain.DeleteMenuAuth(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}

// GetTenantMenu 获取指定租户的菜单范围（带菜单与按钮权限标记）
// GET /api/v1/admin/platform/menu/tenant?tenant_id={id}
func GetTenantMenu(c *gin.Context) {
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
	tree, err := menudomain.GetTenantMenuTree(uint(tenantIDValue))
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取租户菜单范围失败")
		return
	}
	response.ReturnData(c, tree)
}

// UpdateTenantMenu 更新指定租户的菜单范围与按钮权限范围
// PUT /api/v1/admin/platform/menu/tenant
func UpdateTenantMenu(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		response.ReturnError(c, response.PERMISSION_DENIED, "仅平台管理员可以调整租户菜单范围")
		return
	}
	req := &struct {
		TenantID uint   `json:"tenant_id" binding:"required"`
		MenuData string `json:"menu_data" binding:"required"`
	}{}
	if !middleware.CheckParam(req, c) {
		return
	}
	// 全量覆盖：解析为完整菜单树响应结构
	var menuData []commonmenu.MenuResponse
	if err := json.Unmarshal([]byte(req.MenuData), &menuData); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "menu_data 参数错误")
		return
	}
	// 从全量树中直接提取被勾选的菜单与按钮权限
	tree, err := menudomain.UpdateTenantMenuScope(req.TenantID, menuData)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新租户菜单范围失败")
		return
	}
	response.ReturnData(c, tree)
}
