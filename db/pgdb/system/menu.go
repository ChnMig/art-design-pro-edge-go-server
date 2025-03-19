package system

import (
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

// GetMenuList 获取用户菜单列表及权限
func GetMenuList(userID uint) ([]MenuResponse, error) {
	var menus []MenuResponse
	// 获取用户信息及其角色
	var user User
	if err := pgdb.GetClient().Where(&User{Model: gorm.Model{ID: userID}}).First(&user).Error; err != nil {
		return nil, err
	}
	// 获取该角色关联的所有菜单(包括权限)
	var role Role
	if err := pgdb.GetClient().Preload("Menus").
		Preload("MenuPermissions").
		Where("id = ?", user.RoleID).
		First(&role).Error; err != nil {
		return nil, err
	}
	// 找出所有顶级菜单(ParentID = 0)
	var rootMenus []Menu
	if len(role.Menus) > 0 {
		for _, menu := range role.Menus {
			if menu.ParentID == 0 && menu.Status == 1 {
				rootMenus = append(rootMenus, menu)
			}
		}
	}
	// 递归构建菜单树
	for _, rootMenu := range rootMenus {
		menuResp := convertMenuToResponse(rootMenu)
		buildMenuTree(&menuResp, role.Menus, role.MenuPermissions)
		menus = append(menus, menuResp)
	}
	return menus, nil
}

// 将数据库菜单转换为响应结构
func convertMenuToResponse(menu Menu) MenuResponse {
	return MenuResponse{
		ID:        menu.ID,
		Path:      menu.Path,
		Name:      menu.Name,
		Component: menu.Component,
		Meta: MenuMeta{
			Title:             menu.Title,
			Icon:              menu.Icon,
			KeepAlive:         menu.KeepAlive == 1, // 1表示缓存
			ShowBadge:         menu.ShowBadge == 1, // 1表示显示
			ShowTextBadge:     menu.ShowTextBadge,
			IsHide:            menu.IsHide == 1,    // 1表示隐藏
			IsHideTab:         menu.IsHideTab == 1, // 1表示隐藏
			Link:              menu.Link,
			IsIframe:          menu.IsIframe == 1,          // 1表示是iframe
			IsInMainContainer: menu.IsInMainContainer == 1, // 1表示在主容器中
		},
	}
}

// 递归构建菜单树
func buildMenuTree(parent *MenuResponse, allMenus []Menu, allPermissions []MenuPermission) {
	for _, menu := range allMenus {
		if menu.ParentID == parent.ID && menu.Status == 1 {
			child := convertMenuToResponse(menu)
			// 为子菜单添加权限列表
			for _, perm := range allPermissions {
				if perm.MenuID == menu.ID {
					child.Meta.AuthList = append(child.Meta.AuthList, MenuPermissionResp{
						ID:       perm.ID,
						Title:    perm.Title,
						AuthMark: perm.Mark,
					})
				}
			}
			// 递归处理这个子菜单的子菜单
			buildMenuTree(&child, allMenus, allPermissions)
			parent.Children = append(parent.Children, child)
		}
	}
}
