package menu

import (
    "encoding/json"
    "errors"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "api-server/api/middleware"
    "api-server/api/response"
    commonmenu "api-server/common/menu"
    "api-server/db/pgdb/system"
)

// GetMenuList 获取平台菜单，并基于 tenant_id 标记该组织已拥有的菜单与权限
// GET /api/v1/admin/platform/menu?tenant_id={id}
func GetMenuList(c *gin.Context) {
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

    // 查询全部菜单与权限
    menus, allAuths, err := system.GetMenuData()
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
        return
    }
    // 查询该组织已授权的菜单范围（菜单级别）
    scopeIDs, err := system.GetTenantMenuScopeIDs(uint(tenantIDValue))
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "获取菜单范围失败")
        return
    }
    // 组织按钮权限按菜单范围派生：属于范围内菜单的按钮标记为已拥有
    roleAuthIds := make([]uint, 0)
    if len(scopeIDs) > 0 {
        for _, a := range allAuths {
            for _, mid := range scopeIDs {
                if a.MenuID == mid {
                    roleAuthIds = append(roleAuthIds, a.ID)
                    break
                }
            }
        }
    }

    // 构建带权限标记的菜单树（all=true 包含所有菜单）
    menuTree := commonmenu.BuildMenuTreeWithPermission(menus, allAuths, scopeIDs, roleAuthIds, true)
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
	if params.ShowBadge == 0 {
		params.ShowBadge = 2
	}
	var level uint = 1
	if params.ParentID != 0 {
		parentMenu := system.SystemMenu{Model: gorm.Model{ID: params.ParentID}}
		if err := system.GetMenu(&parentMenu); err != nil {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
			return
		}
		if parentMenu.Status != system.StatusEnabled {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
			return
		}
		level = parentMenu.Level + 1
	}
	menu := system.SystemMenu{
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
		Level:         level,
		ParentID:      params.ParentID,
		Sort:          params.Sort,
	}
	if err := system.AddMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单失败")
		return
	}
	response.ReturnData(c, menu)
}

// UpdateMenu 根据入参判断：当包含 tenant_id 和 menu_data 时，按组织范围保存
// PUT /api/v1/admin/platform/menu
func UpdateMenu(c *gin.Context) {
    // 解析组织范围更新请求
    scopeReq := &struct {
        TenantID uint   `json:"tenant_id" form:"tenant_id"`
        MenuData string `json:"menu_data" form:"menu_data"`
    }{}
    // 先不校验，尝试绑定；若 TenantID 存在则按范围更新处理
    _ = c.ShouldBind(scopeReq)
    if scopeReq.TenantID != 0 && scopeReq.MenuData != "" {
        // 鉴权：平台管理员
        if !middleware.IsSuperAdmin(c) {
            response.ReturnError(c, response.PERMISSION_DENIED, "仅平台管理员可以调整租户菜单范围")
            return
        }
        // 反序列化 menu_data
        var menuData []commonmenu.MenuResponse
        if err := json.Unmarshal([]byte(scopeReq.MenuData), &menuData); err != nil {
            response.ReturnError(c, response.INVALID_ARGUMENT, "menu_data 参数错误")
            return
        }
        // 提取被勾选的菜单ID
        menuIDs := extractCheckedMenuIDs(menuData)
        if err := system.SaveTenantMenuScope(scopeReq.TenantID, menuIDs); err != nil {
            response.ReturnError(c, response.DATA_LOSS, "保存菜单范围失败")
            return
        }
        response.ReturnData(c, gin.H{"tenant_id": scopeReq.TenantID, "menu_ids": menuIDs})
        return
    }

    // 否则按“菜单定义更新”处理（与原先一致）
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
    if params.ShowBadge == 0 {
        params.ShowBadge = 2
    }
    var level uint = 1
    if params.ParentID != 0 {
        parent := system.SystemMenu{Model: gorm.Model{ID: params.ParentID}}
        if err := system.GetMenu(&parent); err != nil {
            response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
            return
        }
        if parent.Status != system.StatusEnabled {
            response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
            return
        }
        level = parent.Level + 1
    }
    if params.Status == system.StatusDisabled {
        children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: params.ID}, -1, -1)
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
            return
        }
        for _, child := range children {
            if child.Status == system.StatusEnabled {
                response.ReturnError(c, response.DATA_LOSS, "请先禁用子菜单")
                return
            }
        }
    }
    menu := system.SystemMenu{
        Model:         gorm.Model{ID: params.ID},
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
        Level:         level,
        ParentID:      params.ParentID,
        Sort:          params.Sort,
    }
    if err := system.UpdateMenu(&menu); err != nil {
        response.ReturnError(c, response.DATA_LOSS, "更新菜单失败")
        return
    }
    response.ReturnData(c, menu)
}

func DeleteMenu(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	menu := system.SystemMenu{Model: gorm.Model{ID: params.ID}}
	if err := system.GetMenu(&menu); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "菜单不存在")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: menu.ID}, -1, -1)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
		return
	}
	if len(children) > 0 {
		response.ReturnError(c, response.DATA_LOSS, "请先删除子菜单")
		return
	}
	if err := system.DeleteMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单失败")
		return
	}
	response.ReturnData(c, menu)
}

func GetMenuAuthList(c *gin.Context) {
	params := &struct {
		MenuID uint `json:"menu_id" form:"menu_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.SystemMenuAuth{MenuID: params.MenuID}
	auths, err := system.FindMenuAuthList(&auth)
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
	auth := system.SystemMenuAuth{
		MenuID: params.MenuID,
		Mark:   params.Mark,
		Title:  params.Title,
	}
	if err := system.AddMenuAuth(&auth); err != nil {
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
	auth := system.SystemMenuAuth{
		Model:  gorm.Model{ID: params.ID},
		Title:  params.Title,
		Mark:   params.Mark,
		MenuID: params.MenuID,
	}
	if err := system.UpdateMenuAuth(&auth); err != nil {
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
	auth := system.SystemMenuAuth{Model: gorm.Model{ID: params.ID}}
	if err := system.DeleteMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}

// extractCheckedMenuIDs 递归提取树中被勾选的菜单ID
func extractCheckedMenuIDs(tree []commonmenu.MenuResponse) []uint {
    var result []uint
    var walk func(items []commonmenu.MenuResponse)
    walk = func(items []commonmenu.MenuResponse) {
        for _, m := range items {
            if m.HasPermission {
                result = append(result, m.ID)
            }
            if len(m.Children) > 0 {
                walk(m.Children)
            }
        }
    }
    walk(tree)
    return result
}
