package system

import (
	"go.uber.org/zap"

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
	if err := pgdb.GetClient().Delete(menuAuth).Error; err != nil {
		zap.L().Error("failed to delete menu Auth", zap.Error(err))
		return err
	}
	return nil
}

func AddMenuAuth(menuAuth *SystemMenuAuth) error {
	if err := pgdb.GetClient().Create(menuAuth).Error; err != nil {
		zap.L().Error("failed to create menu Auth", zap.Error(err))
		return err
	}
	return nil
}

func UpdateMenuAuth(menuAuth *SystemMenuAuth) error {
	if err := pgdb.GetClient().Updates(menuAuth).Error; err != nil {
		zap.L().Error("failed to update menu Auth", zap.Error(err))
		return err
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
