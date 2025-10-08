package menu

import (
    "encoding/json"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "api-server/api/middleware"
    "api-server/api/response"
    "api-server/common/menu"
    "api-server/db/pgdb/system"
)

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
        // 先按“菜单范围”过滤可见的菜单与其按钮定义
        allMenus, allAuths = system.FilterMenusByIDs(allMenus, allAuths, scopeIDs)

        // 再按“按钮权限范围”过滤按钮
        authScopeIDs, err := system.GetTenantAuthScopeIDs(roleEntity.TenantID)
        if err != nil {
            response.ReturnError(c, response.DATA_LOSS, "获取按钮权限范围失败")
            return
        }
        if len(authScopeIDs) > 0 {
            allowedAuthSet := make(map[uint]struct{}, len(authScopeIDs))
            for _, id := range authScopeIDs { allowedAuthSet[id] = struct{}{} }
            filteredAuths := make([]system.SystemMenuAuth, 0, len(allAuths))
            for _, a := range allAuths {
                if _, ok := allowedAuthSet[a.ID]; ok { filteredAuths = append(filteredAuths, a) }
            }
            allAuths = filteredAuths
        } else {
            allAuths = []system.SystemMenuAuth{}
        }

        // 过滤角色当前拥有的菜单/按钮集合到允许集合内
        if len(allMenus) == 0 {
            roleMenuIds = []uint{}
            roleAuthIds = []uint{}
        } else {
            allowedMenuIDs := make([]uint, 0, len(allMenus))
            for _, m := range allMenus { allowedMenuIDs = append(allowedMenuIDs, m.ID) }
            roleMenuIds = system.FilterUintIDs(roleMenuIds, allowedMenuIDs)

            allowedAuthIDs := make([]uint, 0, len(allAuths))
            for _, auth := range allAuths { allowedAuthIDs = append(allowedAuthIDs, auth.ID) }
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
