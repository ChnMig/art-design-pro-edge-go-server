package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// SaveRoleMenuAssociations 根据菜单/按钮权限 ID 集合，全量覆盖角色的菜单与按钮权限关联。
// 说明：
// - 菜单（SystemMenus）与按钮权限（SystemMenuAuths）独立保存；
// - 由上层（domain/api）负责完成“可分配范围”校验与 ID 提取。
func SaveRoleMenuAssociations(roleID uint, menuIDs []uint, authIDs []uint) error {
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		var role SystemRole
		if err := tx.First(&role, roleID).Error; err != nil {
			zap.L().Error("failed to find role", zap.Uint("role_id", roleID), zap.Error(err))
			return err
		}

		if err := tx.Model(&role).Association("SystemMenus").Clear(); err != nil {
			zap.L().Error("failed to clear role menus", zap.Uint("role_id", roleID), zap.Error(err))
			return err
		}
		if len(menuIDs) > 0 {
			var menus []SystemMenu
			if err := tx.Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
				zap.L().Error("failed to find menus", zap.Uint("role_id", roleID), zap.Uints("menu_ids", menuIDs), zap.Error(err))
				return err
			}
			if err := tx.Model(&role).Association("SystemMenus").Append(&menus); err != nil {
				zap.L().Error("failed to append menus to role", zap.Uint("role_id", roleID), zap.Error(err))
				return err
			}
		}

		if err := tx.Model(&role).Association("SystemMenuAuths").Clear(); err != nil {
			zap.L().Error("failed to clear role auths", zap.Uint("role_id", roleID), zap.Error(err))
			return err
		}
		if len(authIDs) > 0 {
			var auths []SystemMenuAuth
			if err := tx.Where("id IN ?", authIDs).Find(&auths).Error; err != nil {
				zap.L().Error("failed to find menu auths", zap.Uint("role_id", roleID), zap.Uints("auth_ids", authIDs), zap.Error(err))
				return err
			}
			if err := tx.Model(&role).Association("SystemMenuAuths").Append(&auths); err != nil {
				zap.L().Error("failed to append auths to role", zap.Uint("role_id", roleID), zap.Error(err))
				return err
			}
		}

		return nil
	})
}

