package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

func GetMenuAuth(auth *SystemMenuAuth) error {
	if err := pgdb.GetClient().Where(auth).First(auth).Error; err != nil {
		zap.L().Error("failed to get menu Auth", zap.Error(err))
		return err
	}
	return nil
}

func DeleteMenuAuth(menuAuth *SystemMenuAuth) error {
	if menuAuth.ID == 0 {
		return gorm.ErrInvalidData
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM system_roles__system_auths WHERE system_menu_auth_id = ?", menuAuth.ID).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("auth_id = ?", menuAuth.ID).Delete(&SystemTenantAuthScope{}).Error; err != nil {
			return err
		}
		result := tx.Where("id = ?", menuAuth.ID).Delete(&SystemMenuAuth{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func AddMenuAuth(menuAuth *SystemMenuAuth) error {
	if err := pgdb.GetClient().Create(menuAuth).Error; err != nil {
		zap.L().Error("failed to create menu Auth", zap.Error(err))
		return err
	}
	return nil
}

func UpdateMenuAuth(menuAuth *SystemMenuAuth) error {
	if menuAuth.ID == 0 {
		return gorm.ErrInvalidData
	}
	result := pgdb.GetClient().Model(&SystemMenuAuth{}).Where("id = ?", menuAuth.ID).
		Select("title", "mark", "menu_id").Updates(menuAuth)
	if result.Error != nil {
		zap.L().Error("failed to update menu Auth", zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func FindMenuAuthList(menuAuth *SystemMenuAuth) ([]SystemMenuAuth, error) {
	var auths []SystemMenuAuth
	if err := pgdb.GetClient().Where(menuAuth).Find(&auths).Error; err != nil {
		zap.L().Error("failed to find menu Auth list", zap.Error(err))
		return nil, err
	}
	return auths, nil
}
