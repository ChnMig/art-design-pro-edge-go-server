package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// GetUserMenuData 获取用户菜单数据
func GetUserMenuData(userID uint) ([]Menu, []MenuAuth, error) {
	// 获取用户信息及其角色
	var user User
	if err := pgdb.GetClient().Where(&User{Model: gorm.Model{ID: userID}}).First(&user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return nil, nil, err
	}
	// 获取该角色关联的所有菜单(包括权限)
	var role Role
	if err := pgdb.GetClient().Preload("Menus").
		Preload("MenuAuth").
		Where("id = ?", user.RoleID).
		First(&role).Error; err != nil {
		zap.L().Error("failed to get role", zap.Error(err))
		return nil, nil, err
	}
	return role.Menus, role.MenuAuth, nil
}

// 获取菜单树(不带分页)
func GetMenuData() ([]Menu, []MenuAuth, error) {
	var menus []Menu
	if err := pgdb.GetClient().Find(&menus).Error; err != nil {
		zap.L().Error("failed to get menus", zap.Error(err))
		return nil, nil, err
	}
	var Auths []MenuAuth
	if err := pgdb.GetClient().Find(&Auths).Error; err != nil {
		zap.L().Error("failed to get menu Auths", zap.Error(err))
		return nil, nil, err
	}
	return menus, Auths, nil
}

// 新增一个菜单
func AddMenu(menu *Menu) error {
	if err := pgdb.GetClient().Create(&menu).Error; err != nil {
		zap.L().Error("failed to create menu", zap.Error(err))
		return err
	}
	return nil
}

// 删除一个菜单
func DeleteMenu(menu *Menu) error {
	if err := pgdb.GetClient().Delete(&menu).Error; err != nil {
		zap.L().Error("failed to delete menu", zap.Error(err))
		return err
	}
	return nil
}

func UpdateMenu(menu *Menu) error {
	if err := pgdb.GetClient().Save(&menu).Error; err != nil {
		zap.L().Error("failed to update menu", zap.Error(err))
		return err
	}
	return nil
}

func GetMenu(menu *Menu) error {
	if err := pgdb.GetClient().Where(menu).First(menu).Error; err != nil {
		zap.L().Error("failed to get menu", zap.Error(err))
		return err
	}
	return nil
}

func FindMenuList(menu *Menu) ([]Menu, error) {
	var menus []Menu
	if err := pgdb.GetClient().Where(menu).Find(&menus).Error; err != nil {
		zap.L().Error("failed to find menu list", zap.Error(err))
		return nil, err
	}
	return menus, nil
}
