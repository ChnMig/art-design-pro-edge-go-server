package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

func GetTenantMenuScopeIDs(tenantID uint) ([]uint, error) {
	if tenantID == 0 {
		return nil, nil
	}
	var scopes []SystemTenantMenuScope
	if err := pgdb.GetClient().Where("tenant_id = ?", tenantID).Find(&scopes).Error; err != nil {
		zap.L().Error("failed to get tenant menu scope", zap.Uint("tenantID", tenantID), zap.Error(err))
		return nil, err
	}
	menuIDs := make([]uint, 0, len(scopes))
	for _, scope := range scopes {
		menuIDs = append(menuIDs, scope.MenuID)
	}
	return menuIDs, nil
}

func SaveTenantMenuScope(tenantID uint, menuIDs []uint) error {
	if tenantID == 0 {
		return fmt.Errorf("tenant id is required")
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		return saveTenantMenuScopeTx(tx, tenantID, uniqueUintIDs(menuIDs))
	})
}

func saveTenantMenuScopeTx(tx *gorm.DB, tenantID uint, menuIDs []uint) error {
	if err := tx.Unscoped().Where("tenant_id = ?", tenantID).Delete(&SystemTenantMenuScope{}).Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}
	var count int64
	if err := tx.Model(&SystemMenu{}).Where("id IN ?", menuIDs).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(menuIDs)) {
		return fmt.Errorf("invalid menu ids provided")
	}
	items := make([]SystemTenantMenuScope, len(menuIDs))
	for i, id := range menuIDs {
		items[i] = SystemTenantMenuScope{TenantID: tenantID, MenuID: id}
	}
	return tx.Create(&items).Error
}

// SaveTenantPermissionScope 原子更新租户的菜单、按钮权限范围，并清理角色的越权关联。
func SaveTenantPermissionScope(tenantID uint, menuIDs, authIDs []uint) error {
	if tenantID == 0 {
		return fmt.Errorf("tenant id is required")
	}
	menuIDs = uniqueUintIDs(menuIDs)
	authIDs = uniqueUintIDs(authIDs)
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		var tenant SystemTenant
		if err := tx.First(&tenant, tenantID).Error; err != nil {
			return err
		}
		if err := validateAuthMenuScopeTx(tx, menuIDs, authIDs); err != nil {
			return err
		}
		if err := saveTenantMenuScopeTx(tx, tenantID, menuIDs); err != nil {
			return err
		}
		if err := saveTenantAuthScopeTx(tx, tenantID, authIDs); err != nil {
			return err
		}
		return pruneTenantRoleAssociationsTx(tx, tenantID, menuIDs, authIDs)
	})
}

func validateAuthMenuScopeTx(tx *gorm.DB, menuIDs, authIDs []uint) error {
	if len(authIDs) == 0 {
		return nil
	}
	if len(menuIDs) == 0 {
		return fmt.Errorf("button permissions require selected menus")
	}
	var count int64
	if err := tx.Model(&SystemMenuAuth{}).
		Where("id IN ? AND menu_id IN ?", authIDs, menuIDs).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(authIDs)) {
		return fmt.Errorf("button permission is outside selected menu scope")
	}
	return nil
}

func uniqueUintIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, found := seen[id]; found {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

// PruneTenantRoleAssociations 当平台调整租户的菜单/按钮范围后，
// 自动清理该租户下所有角色中超出范围的角色-菜单/角色-按钮关联。
func PruneTenantRoleAssociations(tenantID uint, allowedMenuIDs []uint, allowedAuthIDs []uint) error {
	if tenantID == 0 {
		return fmt.Errorf("tenant id is required")
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		return pruneTenantRoleAssociationsTx(tx, tenantID, allowedMenuIDs, allowedAuthIDs)
	})
}

func pruneTenantRoleAssociationsTx(tx *gorm.DB, tenantID uint, allowedMenuIDs []uint, allowedAuthIDs []uint) error {
	// 查找该租户下所有角色ID
	var roleIDs []uint
	if err := tx.Model(&SystemRole{}).Where("tenant_id = ?", tenantID).Pluck("id", &roleIDs).Error; err != nil {
		zap.L().Error("failed to list tenant roles", zap.Uint("tenantID", tenantID), zap.Error(err))
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}

	// 1) 角色-菜单：删除超出 allowedMenuIDs 的关联
	if len(allowedMenuIDs) > 0 {
		if err := tx.Exec(
			"DELETE FROM system_roles__system_menus WHERE system_role_id IN ? AND system_menu_id NOT IN ?",
			roleIDs, allowedMenuIDs,
		).Error; err != nil {
			zap.L().Error("failed to prune role-menu associations", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
	} else {
		// 未配置菜单范围，清空所有角色-菜单关联
		if err := tx.Exec(
			"DELETE FROM system_roles__system_menus WHERE system_role_id IN ?",
			roleIDs,
		).Error; err != nil {
			zap.L().Error("failed to clear role-menu associations", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
	}

	// 2) 角色-按钮：删除超出 allowedAuthIDs 的关联
	if len(allowedAuthIDs) > 0 {
		if err := tx.Exec(
			"DELETE FROM system_roles__system_auths WHERE system_role_id IN ? AND system_menu_auth_id NOT IN ?",
			roleIDs, allowedAuthIDs,
		).Error; err != nil {
			zap.L().Error("failed to prune role-auth associations", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
	} else {
		// 未配置按钮范围，清空所有角色-按钮关联
		if err := tx.Exec(
			"DELETE FROM system_roles__system_auths WHERE system_role_id IN ?",
			roleIDs,
		).Error; err != nil {
			zap.L().Error("failed to clear role-auth associations", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
	}

	// 3) 保护性清理：若某些按钮所属菜单不在 allowedMenuIDs 内，一并移除按钮关联
	if len(allowedMenuIDs) > 0 {
		if err := tx.Exec(
			"DELETE FROM system_roles__system_auths ra USING system_menu_auths a "+
				"WHERE ra.system_menu_auth_id = a.id AND ra.system_role_id IN ? AND a.menu_id NOT IN ?",
			roleIDs, allowedMenuIDs,
		).Error; err != nil {
			zap.L().Error("failed to prune role-auth by disallowed menus", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
	}

	return nil
}
