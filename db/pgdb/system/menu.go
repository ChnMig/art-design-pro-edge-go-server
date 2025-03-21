package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// MenuResponse 定义返回给前端的菜单结构
type MenuResponse struct {
	ID        uint           `json:"id"`
	Path      string         `json:"path"`
	Name      string         `json:"name"`
	Component string         `json:"component,omitempty"`
	Meta      MenuMeta       `json:"meta"`
	Children  []MenuResponse `json:"children,omitempty"`
}

// MenuMeta 定义菜单元数据
type MenuMeta struct {
	Title             string               `json:"title"`
	Icon              string               `json:"icon,omitempty"`
	KeepAlive         bool                 `json:"keepAlive"`
	ShowBadge         bool                 `json:"showBadge,omitempty"`
	ShowTextBadge     string               `json:"showTextBadge,omitempty"`
	IsHide            bool                 `json:"isHide,omitempty"`
	IsHideTab         bool                 `json:"isHideTab,omitempty"`
	Link              string               `json:"link,omitempty"`
	IsIframe          bool                 `json:"isIframe,omitempty"`
	IsInMainContainer bool                 `json:"isInMainContainer,omitempty"`
	AuthList          []MenuPermissionResp `json:"authList,omitempty"`
}

// MenuPermissionResp 定义菜单权限响应结构
type MenuPermissionResp struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	AuthMark string `json:"auth_mark"` // 对应 MenuPermission 中的 Mark
}

// GetUserMenuData 获取用户菜单数据
func GetUserMenuData(userID uint) ([]Menu, []MenuPermission, error) {
	// 获取用户信息及其角色
	var user User
	if err := pgdb.GetClient().Where(&User{Model: gorm.Model{ID: userID}}).First(&user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return nil, nil, err
	}
	// 获取该角色关联的所有菜单(包括权限)
	var role Role
	if err := pgdb.GetClient().Preload("Menus").
		Preload("MenuPermissions").
		Where("id = ?", user.RoleID).
		First(&role).Error; err != nil {
		zap.L().Error("failed to get role", zap.Error(err))
		return nil, nil, err
	}
	return role.Menus, role.MenuPermissions, nil
}

// 获取菜单树(不带分页)
func GetMenuData() ([]Menu, []MenuPermission, error) {
	var menus []Menu
	if err := pgdb.GetClient().Find(&menus).Error; err != nil {
		zap.L().Error("failed to get menus", zap.Error(err))
		return nil, nil, err
	}
	var permissions []MenuPermission
	if err := pgdb.GetClient().Find(&permissions).Error; err != nil {
		zap.L().Error("failed to get menu permissions", zap.Error(err))
		return nil, nil, err
	}
	return menus, permissions, nil
}

// 新增一个菜单
func AddMenu(menu *Menu) error {
	if err := pgdb.GetClient().Create(&menu).Error; err != nil {
		zap.L().Error("failed to create menu", zap.Error(err))
		return err
	}
	return nil
}

// 通过 ID 获取菜单
func GetMenuByID(id uint) (Menu, error) {
	var menu Menu
	if err := pgdb.GetClient().Where("id = ?", id).First(&menu).Error; err != nil {
		zap.L().Error("failed to get menu by id", zap.Error(err))
		return menu, err
	}
	return menu, nil
}
