package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// GetUserMenuData 获取用户菜单数据
func GetUserMenuData(userID uint) ([]SystemMenu, []SystemMenuAuth, error) {
	// 获取用户信息及其角色
	var user SystemUser
	if err := pgdb.GetClient().Where(&SystemUser{Model: gorm.Model{ID: userID}}).First(&user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return nil, nil, err
	}
	// 获取该角色关联的所有菜单(包括权限)
	var role SystemRole
	if err := pgdb.GetClient().Preload("SystemMenus").
		Preload("SystemMenuAuths").
		Where("id = ?", user.RoleID).
		First(&role).Error; err != nil {
		zap.L().Error("failed to get role", zap.Error(err))
		return nil, nil, err
	}
	return role.SystemMenus, role.SystemMenuAuths, nil
}

// 获取菜单树(不带分页)
func GetMenuData() ([]SystemMenu, []SystemMenuAuth, error) {
	var menus []SystemMenu
	if err := pgdb.GetClient().Find(&menus).Error; err != nil {
		zap.L().Error("failed to get menus", zap.Error(err))
		return nil, nil, err
	}
	var Auths []SystemMenuAuth
	if err := pgdb.GetClient().Find(&Auths).Error; err != nil {
		zap.L().Error("failed to get menu Auths", zap.Error(err))
		return nil, nil, err
	}
	return menus, Auths, nil
}

// GetMenuDataByRoleID 获取指定角色ID的菜单和权限数据
func GetMenuDataByRoleID(roleID uint) ([]SystemMenu, []SystemMenuAuth, []uint, []uint, error) {
	// 获取所有菜单
	var allMenus []SystemMenu
	if err := pgdb.GetClient().Find(&allMenus).Error; err != nil {
		zap.L().Error("failed to get all menus", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	// 获取所有权限
	var allAuths []SystemMenuAuth
	if err := pgdb.GetClient().Find(&allAuths).Error; err != nil {
		zap.L().Error("failed to get all menu auths", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	// 获取角色拥有的菜单ID列表
	var role SystemRole
	if err := pgdb.GetClient().Preload("SystemMenus").
		Preload("SystemMenuAuths").
		Where("id = ?", roleID).
		First(&role).Error; err != nil {
		zap.L().Error("failed to get role with menus", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	// 提取角色拥有的菜单ID和权限ID
	var roleMenuIds []uint
	var roleAuthIds []uint
	for _, m := range role.SystemMenus {
		roleMenuIds = append(roleMenuIds, m.ID)
	}
	for _, a := range role.SystemMenuAuths {
		roleAuthIds = append(roleAuthIds, a.ID)
	}
	return allMenus, allAuths, roleMenuIds, roleAuthIds, nil
}

// 新增一个菜单
func AddMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Create(&menu).Error; err != nil {
		zap.L().Error("failed to create menu", zap.Error(err))
		return err
	}
	return nil
}

// 删除一个菜单
func DeleteMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Delete(&menu).Error; err != nil {
		zap.L().Error("failed to delete menu", zap.Error(err))
		return err
	}
	return nil
}

func UpdateMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Save(&menu).Error; err != nil {
		zap.L().Error("failed to update menu", zap.Error(err))
		return err
	}
	return nil
}

func GetMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Where(menu).First(menu).Error; err != nil {
		zap.L().Error("failed to get menu", zap.Error(err))
		return err
	}
	return nil
}

func FindMenuList(menu *SystemMenu) ([]SystemMenu, error) {
	var menus []SystemMenu
	if err := pgdb.GetClient().Where(menu).Find(&menus).Error; err != nil {
		zap.L().Error("failed to find menu list", zap.Error(err))
		return nil, err
	}
	return menus, nil
}
