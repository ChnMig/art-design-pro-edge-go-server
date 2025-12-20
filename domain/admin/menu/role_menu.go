package menu

import (
	"errors"

	commonmenu "api-server/common/menu"
	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

func GetRoleMenuTree(roleID uint, actorTenantID uint, isSuperAdmin bool) ([]commonmenu.MenuResponse, error) {
	roleEntity := system.SystemRole{Model: gorm.Model{ID: roleID}}
	if err := system.GetRole(&roleEntity); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	if !isSuperAdmin {
		if actorTenantID == 0 || roleEntity.TenantID != actorTenantID {
			return nil, ErrPermissionDenied
		}
	}

	allMenus, allAuths, roleMenuIDs, roleAuthIDs, err := system.GetMenuDataByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	scopeIDs, err := system.GetTenantMenuScopeIDs(roleEntity.TenantID)
	if err != nil {
		return nil, err
	}
	allMenus, allAuths = system.FilterMenusByIDs(allMenus, allAuths, scopeIDs)

	authScopeIDs, err := system.GetTenantAuthScopeIDs(roleEntity.TenantID)
	if err != nil {
		return nil, err
	}
	allAuths = filterAuthsByScope(allAuths, authScopeIDs)

	if len(allMenus) == 0 {
		roleMenuIDs = []uint{}
		roleAuthIDs = []uint{}
	} else {
		allowedMenuIDs := make([]uint, 0, len(allMenus))
		for _, m := range allMenus {
			allowedMenuIDs = append(allowedMenuIDs, m.ID)
		}
		roleMenuIDs = system.FilterUintIDs(roleMenuIDs, allowedMenuIDs)

		allowedAuthIDs := make([]uint, 0, len(allAuths))
		for _, auth := range allAuths {
			allowedAuthIDs = append(allowedAuthIDs, auth.ID)
		}
		roleAuthIDs = system.FilterUintIDs(roleAuthIDs, allowedAuthIDs)
	}

	return commonmenu.BuildMenuTreeWithPermission(allMenus, allAuths, roleMenuIDs, roleAuthIDs, true), nil
}

func UpdateRoleMenu(roleID uint, menuData []commonmenu.MenuResponse, actorTenantID uint, isSuperAdmin bool) error {
	roleEntity := system.SystemRole{Model: gorm.Model{ID: roleID}}
	if err := system.GetRole(&roleEntity); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	if !isSuperAdmin {
		if actorTenantID == 0 || actorTenantID != roleEntity.TenantID {
			return ErrPermissionDenied
		}
	}

	scopeIDs, err := system.GetTenantMenuScopeIDs(roleEntity.TenantID)
	if err != nil {
		return err
	}
	if !validateMenuScope(menuData, scopeIDs) {
		return ErrMenuOutOfScope
	}

	authScopeIDs, err := system.GetTenantAuthScopeIDs(roleEntity.TenantID)
	if err != nil {
		return err
	}
	if !validateAuthScope(menuData, authScopeIDs) {
		return ErrAuthOutOfScope
	}

	menuIDs := extractCheckedMenuIDs(menuData)
	authIDs := extractCheckedAuthIDs(menuData)

	return system.SaveRoleMenuAssociations(roleID, menuIDs, authIDs)
}

