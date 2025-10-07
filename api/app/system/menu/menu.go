package menu

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/common/menu"
	"api-server/db/pgdb/system"
)

func GetMenuList(c *gin.Context) {
	isSuperAdmin := middleware.IsSuperAdmin(c)
	tenantID := middleware.GetTenantID(c)

	// 查询菜单数据
	menus, menup, err := system.GetMenuData()
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	if !isSuperAdmin {
		if tenantID == 0 {
			response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
			return
		}
		scopeIDs, err := system.GetTenantMenuScopeIDs(tenantID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取菜单范围失败")
			return
		}
		menus, menup = system.FilterMenusByIDs(menus, menup, scopeIDs)
	}
	// 构建菜单树
	menuTree := menu.BuildMenuTree(menus, menup, true)
	response.ReturnData(c, menuTree)
}

func DeleteMenu(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		response.ReturnError(c, response.PERMISSION_DENIED, "仅平台管理员可以删除菜单")
		return
	}
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	menu := system.SystemMenu{Model: gorm.Model{ID: params.ID}}
	err := system.GetMenu(&menu)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "菜单不存在")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: menu.ID}, -1, -1)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
			return
		}
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

func AddMenu(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		response.ReturnError(c, response.PERMISSION_DENIED, "仅平台管理员可以创建菜单")
		return
	}
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
	// 如果有父级ID, 则查询父级ID是否存在
	if params.ParentID != 0 {
		parentMenu := system.SystemMenu{
			Model: gorm.Model{ID: params.ParentID},
		}
		err := system.GetMenu(&parentMenu)
		if err != nil {
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

func UpdateMenu(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		response.ReturnError(c, response.PERMISSION_DENIED, "仅平台管理员可以修改菜单")
		return
	}
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
	// 如果有父级ID, 则查询父级ID是否存在
	if params.ParentID != 0 {
		parent := system.SystemMenu{Model: gorm.Model{ID: params.ParentID}}
		err := system.GetMenu(&parent)
		if err != nil {
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
		// 判断子菜单是否是禁用状态
		children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: params.ID}, -1, -1)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
				return
			}
		}
		if len(children) > 0 {
			for _, v := range children {
				if v.Status == system.StatusEnabled {
					response.ReturnError(c, response.DATA_LOSS, "请先禁用子菜单")
					return
				}
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

// GetMenuListByRoleID 根据角色ID获取菜单列表
func GetMenuListByRoleID(c *gin.Context) {
	params := &struct {
		RoleID uint `json:"role_id" form:"role_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	isSuperAdmin := middleware.IsSuperAdmin(c)
	currentTenantID := middleware.GetTenantID(c)

	roleEntity := system.SystemRole{Model: gorm.Model{ID: params.RoleID}}
	if err := system.GetRole(&roleEntity); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "角色不存在")
		return
	}
	if !isSuperAdmin {
		if currentTenantID == 0 || roleEntity.TenantID != currentTenantID {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权查看该角色菜单")
			return
		}
	}
	// 查询菜单数据
	allMenus, allAuths, roleMenuIds, roleAuthIds, err := system.GetMenuDataByRoleID(params.RoleID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询角色菜单失败")
		return
	}
	if !isSuperAdmin {
		scopeIDs, err := system.GetTenantMenuScopeIDs(roleEntity.TenantID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取菜单范围失败")
			return
		}
		allMenus, allAuths = system.FilterMenusByIDs(allMenus, allAuths, scopeIDs)
		if len(allMenus) == 0 {
			roleMenuIds = []uint{}
			roleAuthIds = []uint{}
		} else {
			allowedMenuIDs := make([]uint, 0, len(allMenus))
			for _, m := range allMenus {
				allowedMenuIDs = append(allowedMenuIDs, m.ID)
			}
			roleMenuIds = system.FilterUintIDs(roleMenuIds, allowedMenuIDs)
			allowedAuthIDs := make([]uint, 0, len(allAuths))
			for _, auth := range allAuths {
				allowedAuthIDs = append(allowedAuthIDs, auth.ID)
			}
			roleAuthIds = system.FilterUintIDs(roleAuthIds, allowedAuthIDs)
		}
	}
	// 构建带权限标记的菜单树
	menuTree := menu.BuildMenuTreeWithPermission(allMenus, allAuths, roleMenuIds, roleAuthIds, true)
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
	isSuperAdmin := middleware.IsSuperAdmin(c)
	// 尝试将 params.MenuData 转成结构体
	var menuData []menu.MenuResponse
	err := json.Unmarshal([]byte(params.MenuData), &menuData)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "参数错误")
		return
	}

	roleEntity := system.SystemRole{Model: gorm.Model{ID: params.RoleID}}
	if err := system.GetRole(&roleEntity); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "角色不存在")
		return
	}
	if !isSuperAdmin {
		tenantID := middleware.GetTenantID(c)
		if tenantID == 0 || tenantID != roleEntity.TenantID {
			response.ReturnError(c, response.PERMISSION_DENIED, "无权调整该角色菜单")
			return
		}
		scopeIDs, err := system.GetTenantMenuScopeIDs(roleEntity.TenantID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取菜单范围失败")
			return
		}
		if !validateMenuScope(menuData, scopeIDs) {
			response.ReturnError(c, response.PERMISSION_DENIED, "菜单超出可分配范围")
			return
		}
	}
	// 保存角色菜单数据
	err = menu.SaveRoleMenu(params.RoleID, menuData)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "保存角色菜单失败")
		return
	}

	response.ReturnData(c, nil)
}

func validateMenuScope(menus []menu.MenuResponse, allowed []uint) bool {
	if len(allowed) == 0 {
		return len(menus) == 0
	}
	allowedSet := make(map[uint]struct{}, len(allowed))
	for _, id := range allowed {
		allowedSet[id] = struct{}{}
	}
	var walk func(items []menu.MenuResponse) bool
	walk = func(items []menu.MenuResponse) bool {
		for _, m := range items {
			if _, ok := allowedSet[m.ID]; !ok {
				return false
			}
			if len(m.Children) > 0 && !walk(m.Children) {
				return false
			}
		}
		return true
	}
	return walk(menus)
}
