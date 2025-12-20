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
		if err := tx.Where("tenant_id = ?", tenantID).Delete(&SystemTenantMenuScope{}).Error; err != nil {
			zap.L().Error("failed to clear tenant menu scope", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		var count int64
		if err := tx.Model(&SystemMenu{}).Where("id IN ?", menuIDs).Count(&count).Error; err != nil {
			zap.L().Error("failed to validate menu ids", zap.Error(err))
			return err
		}
		if count != int64(len(menuIDs)) {
			return fmt.Errorf("invalid menu ids provided")
		}
		items := make([]SystemTenantMenuScope, len(menuIDs))
		for i, id := range menuIDs {
			items[i] = SystemTenantMenuScope{
				TenantID: tenantID,
				MenuID:   id,
			}
		}
		if err := tx.Create(&items).Error; err != nil {
			zap.L().Error("failed to create tenant menu scope", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
		return nil
	})
}

// PruneTenantRoleAssociations 当平台调整租户的菜单/按钮范围后，
// 自动清理该租户下所有角色中超出范围的角色-菜单/角色-按钮关联。
func PruneTenantRoleAssociations(tenantID uint, allowedMenuIDs []uint, allowedAuthIDs []uint) error {
	if tenantID == 0 {
		return fmt.Errorf("tenant id is required")
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
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
	})
}
