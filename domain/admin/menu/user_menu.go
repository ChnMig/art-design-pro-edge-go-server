package menu

import (
	commonmenu "api-server/common/menu"
	"api-server/db/pgdb/system"
)

func GetUserMenuTree(userID uint, tenantID uint) ([]commonmenu.MenuResponse, error) {
	roleMenus, rolePermissions, err := system.GetUserMenuData(userID)
	if err != nil {
		return nil, err
	}

	scopeIDs, err := system.GetTenantMenuScopeIDs(tenantID)
	if err != nil {
		return nil, err
	}
	roleMenus, rolePermissions = system.FilterMenusByIDs(roleMenus, rolePermissions, scopeIDs)

	authScopeIDs, err := system.GetTenantAuthScopeIDs(tenantID)
	if err != nil {
		return nil, err
	}
	rolePermissions = filterAuthsByScope(rolePermissions, authScopeIDs)

	var roleMenuIDs []uint
	for _, m := range roleMenus {
		roleMenuIDs = append(roleMenuIDs, m.ID)
	}
	var roleAuthIDs []uint
	for _, a := range rolePermissions {
		roleAuthIDs = append(roleAuthIDs, a.ID)
	}

	return commonmenu.BuildMenuTreeWithPermission(roleMenus, rolePermissions, roleMenuIDs, roleAuthIDs, false), nil
}

