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

// GetMenuList 获取平台菜单定义（不带租户 hasPermission 标记）。
// GET /api/v1/admin/platform/menu
func GetMenuList(c *gin.Context) {
    menus, allAuths, err := system.GetMenuData()
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
        return
    }
    menuTree := commonmenu.BuildMenuTree(menus, allAuths, true)
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
    menus, allAuths, err := system.GetMenuData()
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
        return
    }
    scopeIDs, err := system.GetTenantMenuScopeIDs(uint(tenantIDValue))
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "获取菜单范围失败")
        return
    }
    roleAuthIds, err := system.GetTenantAuthScopeIDs(uint(tenantIDValue))
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "获取按钮权限范围失败")
        return
    }
    tree := commonmenu.BuildMenuTreeWithPermission(menus, allAuths, scopeIDs, roleAuthIds, true)
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
    var menuData []commonmenu.MenuResponse
    if err := json.Unmarshal([]byte(req.MenuData), &menuData); err != nil {
        response.ReturnError(c, response.INVALID_ARGUMENT, "menu_data 参数错误")
        return
    }
    menuIDs := extractCheckedMenuIDs(menuData)
    authIDs := extractCheckedAuthIDs(menuData)
    if len(authIDs) > 0 && len(menuIDs) > 0 {
        _, allAuths, err := system.GetMenuData()
        if err != nil {
            response.ReturnError(c, response.DATA_LOSS, "查询菜单权限失败")
            return
        }
        authIDs = filterAuthIDsByMenus(authIDs, menuIDs, allAuths)
    }
    if err := system.SaveTenantMenuScope(req.TenantID, menuIDs); err != nil {
        response.ReturnError(c, response.DATA_LOSS, "保存菜单范围失败")
        return
    }
    if err := system.SaveTenantAuthScope(req.TenantID, authIDs); err != nil {
        response.ReturnError(c, response.DATA_LOSS, "保存按钮权限范围失败")
        return
    }
    menus, allAuths, err := system.GetMenuData()
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
        return
    }
    tree := commonmenu.BuildMenuTreeWithPermission(menus, allAuths, menuIDs, authIDs, true)
    response.ReturnData(c, tree)
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

// extractCheckedAuthIDs 递归提取树中被勾选的按钮权限ID
func extractCheckedAuthIDs(tree []commonmenu.MenuResponse) []uint {
    var result []uint
    var walk func(items []commonmenu.MenuResponse)
    walk = func(items []commonmenu.MenuResponse) {
        for _, m := range items {
            if len(m.Meta.AuthList) > 0 {
                for _, a := range m.Meta.AuthList {
                    if a.HasPermission {
                        result = append(result, a.ID)
                    }
                }
            }
            if len(m.Children) > 0 {
                walk(m.Children)
            }
        }
    }
    walk(tree)
    return result
}

// filterAuthIDsByMenus 过滤按钮权限，仅保留属于提供菜单集合的权限
func filterAuthIDsByMenus(authIDs []uint, menuIDs []uint, allAuths []system.SystemMenuAuth) []uint {
    if len(authIDs) == 0 || len(menuIDs) == 0 {
        return []uint{}
    }
    menuSet := make(map[uint]struct{}, len(menuIDs))
    for _, mid := range menuIDs {
        menuSet[mid] = struct{}{}
    }
    // 建立 authID -> menuID 映射
    authToMenu := make(map[uint]uint, len(allAuths))
    for _, a := range allAuths {
        authToMenu[a.ID] = a.MenuID
    }
    out := make([]uint, 0, len(authIDs))
    for _, aid := range authIDs {
        if mid, ok := authToMenu[aid]; ok {
            if _, allowed := menuSet[mid]; allowed {
                out = append(out, aid)
            }
        }
    }
    return out
}
