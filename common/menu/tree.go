package menu

import (
	"sort"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
	"api-server/db/pgdb/system"
)

// MenuResponse 定义返回给前端的菜单结构
type MenuResponse struct {
	ID            uint           `json:"id"`
	UpdatedAt     uint           `json:"updatedAt,omitempty"`
	Path          string         `json:"path"`
	Name          string         `json:"name"`
	Component     string         `json:"component,omitempty"`
	Meta          MenuMeta       `json:"meta"`
	Children      []MenuResponse `json:"children,omitempty"`
	ParentID      uint           `json:"parentId,omitempty"`
	HasPermission bool           `json:"hasPermission"` // 角色是否拥有此菜单权限
}

// MenuMeta 定义菜单元数据
type MenuMeta struct {
	Title         string         `json:"title"`
	Icon          string         `json:"icon,omitempty"`
	KeepAlive     bool           `json:"keepAlive"`
	ShowBadge     bool           `json:"showBadge,omitempty"`
	ShowTextBadge string         `json:"showTextBadge,omitempty"`
	IsHide        bool           `json:"isHide,omitempty"`
	IsHideTab     bool           `json:"isHideTab,omitempty"`
	Link          string         `json:"link,omitempty"`
	IsIframe      bool           `json:"isIframe,omitempty"`
	IsFirstLevel  bool           `json:"isFirstLevel,omitempty"`
	IsEnable      bool           `json:"isEnable,omitempty"`
	Sort          uint           `json:"sort,omitempty"`
	AuthList      []MenuAuthResp `json:"authList,omitempty"`
}

// 定义菜单权限响应结构
type MenuAuthResp struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	AuthMark      string `json:"authMark"`      // 权限标识
	HasPermission bool   `json:"hasPermission"` // 角色是否拥有此权限
}

// 构建菜单树
// all: 是否包含所有菜单，true表示包含所有菜单，false表示只包含启用的菜单
func BuildMenuTree(menus []system.SystemMenu, permissions []system.SystemMenuAuth, all bool) []MenuResponse {
	var menuTree []MenuResponse

	// 找出所有顶级菜单(ParentID = 0)
	var rootMenus []system.SystemMenu
	if len(menus) > 0 {
		for _, menu := range menus {
			if menu.ParentID == 0 {
				if !all {
					if menu.Status == 2 {
						continue
					}
				}
				rootMenus = append(rootMenus, menu)
			}
		}
	}

	// 对根菜单按 Sort 从大到小排序，Sort 为 0 则按 ID 从大到小
	sortMenus(rootMenus)

	// 递归构建菜单树
	for _, rootMenu := range rootMenus {
		menuResp := convertMenuToResponse(rootMenu)
		buildMenuChildren(&menuResp, menus, permissions, all)
		menuTree = append(menuTree, menuResp)
	}

	return menuTree
}

// 构建带有权限标记的菜单树
func BuildMenuTreeWithPermission(menus []system.SystemMenu, permissions []system.SystemMenuAuth, roleMenuIds []uint, roleAuthIds []uint, all bool) []MenuResponse {
	var menuTree []MenuResponse

	// 找出所有顶级菜单(ParentID = 0)
	var rootMenus []system.SystemMenu
	if len(menus) > 0 {
		for _, menu := range menus {
			if menu.ParentID == 0 {
				if !all && menu.Status == 2 {
					continue
				}
				rootMenus = append(rootMenus, menu)
			}
		}
	}

	// 对根菜单按 Sort 从大到小排序
	sortMenus(rootMenus)

	// 递归构建菜单树
	for _, rootMenu := range rootMenus {
		menuResp := convertMenuToResponseWithPermission(rootMenu, roleMenuIds)
		buildMenuChildrenWithPermission(&menuResp, menus, permissions, roleMenuIds, roleAuthIds, all)
		menuTree = append(menuTree, menuResp)
	}

	return menuTree
}

// 将数据库菜单转换为响应结构
func convertMenuToResponse(menu system.SystemMenu) MenuResponse {
	return MenuResponse{
		ID:        menu.ID,
		UpdatedAt: uint(menu.UpdatedAt.Unix()),
		Path:      menu.Path,
		Name:      menu.Name,
		Component: menu.Component,
		ParentID:  menu.ParentID,
		Meta: MenuMeta{
			Title:         menu.Title,
			Icon:          menu.Icon,
			KeepAlive:     menu.KeepAlive == 1, // 1表示缓存
			ShowBadge:     menu.ShowBadge == 1, // 1表示显示
			ShowTextBadge: menu.ShowTextBadge,
			IsHide:        menu.IsHide == 1,    // 1表示隐藏
			IsHideTab:     menu.IsHideTab == 1, // 1表示隐藏
			Link:          menu.Link,
			IsIframe:      menu.IsIframe == 1,     // 1表示是iframe
			IsFirstLevel:  menu.IsFirstLevel == 1, // 1表示在主容器中
			IsEnable:      menu.Status == 1,       // 1表示启用
			Sort:          menu.Sort,
		},
	}
}

// 将数据库菜单转换为响应结构，并标记是否有权限
func convertMenuToResponseWithPermission(menu system.SystemMenu, roleMenuIds []uint) MenuResponse {
	resp := convertMenuToResponse(menu)
	resp.HasPermission = containsUint(roleMenuIds, menu.ID)
	return resp
}

// 递归构建菜单子项
func buildMenuChildren(parent *MenuResponse, allMenus []system.SystemMenu, allPermissions []system.SystemMenuAuth, all bool) {
	var childMenus []system.SystemMenu

	// 收集当前父菜单下的所有子菜单
	for _, menu := range allMenus {
		if menu.ParentID == parent.ID {
			if !all && menu.Status == 2 {
				continue
			}
			childMenus = append(childMenus, menu)
		}
	}

	// 对子菜单进行排序
	sortMenus(childMenus)

	// 处理排序后的子菜单
	for _, menu := range childMenus {
		child := convertMenuToResponse(menu)

		// 为子菜单添加权限列表
		for _, perm := range allPermissions {
			if perm.MenuID == menu.ID {
				child.Meta.AuthList = append(child.Meta.AuthList, MenuAuthResp{
					ID:       perm.ID,
					Title:    perm.Title,
					AuthMark: perm.Mark,
				})
			}
		}

		// 递归处理这个子菜单的子菜单
		buildMenuChildren(&child, allMenus, allPermissions, all)
		parent.Children = append(parent.Children, child)
	}
}

// 递归构建带权限标记的菜单子项
func buildMenuChildrenWithPermission(parent *MenuResponse, allMenus []system.SystemMenu, allPermissions []system.SystemMenuAuth, roleMenuIds []uint, roleAuthIds []uint, all bool) {
	var childMenus []system.SystemMenu

	// 收集当前父菜单下的所有子菜单
	for _, menu := range allMenus {
		if menu.ParentID == parent.ID {
			if !all && menu.Status == 2 {
				continue
			}
			childMenus = append(childMenus, menu)
		}
	}

	// 对子菜单进行排序
	sortMenus(childMenus)

	// 处理排序后的子菜单
	for _, menu := range childMenus {
		child := convertMenuToResponseWithPermission(menu, roleMenuIds)

		// 为子菜单添加权限列表
		for _, perm := range allPermissions {
			if perm.MenuID == menu.ID {
				authResp := MenuAuthResp{
					ID:            perm.ID,
					Title:         perm.Title,
					AuthMark:      perm.Mark,
					HasPermission: containsUint(roleAuthIds, perm.ID),
				}
				child.Meta.AuthList = append(child.Meta.AuthList, authResp)
			}
		}

		// 递归处理这个子菜单的子菜单
		buildMenuChildrenWithPermission(&child, allMenus, allPermissions, roleMenuIds, roleAuthIds, all)
		parent.Children = append(parent.Children, child)
	}
}

// 对菜单切片按 Sort 从小到大排序，Sort 为 0 则按 ID 从小到大排序（0 视为最大值）
func sortMenus(menus []system.SystemMenu) {
	sort.Slice(menus, func(i, j int) bool {
		// 如果两个菜单都有 Sort 值并且不为 0
		if menus[i].Sort > 0 && menus[j].Sort > 0 {
			// Sort 值不同时，按 Sort 从小到大排序
			if menus[i].Sort != menus[j].Sort {
				return menus[i].Sort < menus[j].Sort
			}
			// Sort 值相同时，按 ID 从小到大排序
			return menus[i].ID < menus[j].ID
		}

		// 如果只有其中一个有 Sort 值
		if menus[i].Sort > 0 && menus[j].Sort == 0 {
			return true // i 排在前面
		}

		if menus[i].Sort == 0 && menus[j].Sort > 0 {
			return false // j 排在前面
		}

		// 都没有 Sort 值或都为 0，按 ID 从小到大排序
		return menus[i].ID < menus[j].ID
	})
}

// 辅助函数：检查uint数组中是否包含特定值
func containsUint(slice []uint, item uint) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func SaveRoleMenu(roleID uint, menuTree []MenuResponse) error {
	// 从menuTree中提取所有有权限的菜单ID和权限ID
	var menuIDs []uint
	var authIDs []uint

	// 递归提取所有有权限的菜单ID和权限ID
	extractPermissions(menuTree, &menuIDs, &authIDs)

	// 使用事务更新数据库
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		// 查询角色
		var role system.SystemRole
		if err := tx.First(&role, roleID).Error; err != nil {
			zap.L().Error("failed to find role", zap.Uint("roleID", roleID), zap.Error(err))
			return err
		}

		// 更新角色的菜单关联
		if err := tx.Model(&role).Association("SystemMenus").Clear(); err != nil {
			zap.L().Error("failed to clear role menus", zap.Uint("roleID", roleID), zap.Error(err))
			return err
		}

		if len(menuIDs) > 0 {
			var menus []system.SystemMenu
			if err := tx.Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
				zap.L().Error("failed to find menus", zap.Uint("roleID", roleID), zap.Uints("menuIDs", menuIDs), zap.Error(err))
				return err
			}
			if err := tx.Model(&role).Association("SystemMenus").Append(&menus); err != nil {
				zap.L().Error("failed to append menus to role", zap.Uint("roleID", roleID), zap.Error(err))
				return err
			}
		}

		// 更新角色的权限关联
		if err := tx.Model(&role).Association("SystemMenuAuths").Clear(); err != nil {
			zap.L().Error("failed to clear role auths", zap.Uint("roleID", roleID), zap.Error(err))
			return err
		}

		if len(authIDs) > 0 {
			var auths []system.SystemMenuAuth
			if err := tx.Where("id IN ?", authIDs).Find(&auths).Error; err != nil {
				zap.L().Error("failed to find menu auths", zap.Uint("roleID", roleID), zap.Uints("authIDs", authIDs), zap.Error(err))
				return err
			}
			if err := tx.Model(&role).Association("SystemMenuAuths").Append(&auths); err != nil {
				zap.L().Error("failed to append auths to role", zap.Uint("roleID", roleID), zap.Error(err))
				return err
			}
		}

		return nil
	})
}

// 递归提取有权限的菜单ID和权限ID
func extractPermissions(menuTree []MenuResponse, menuIDs *[]uint, authIDs *[]uint) {
	for _, menu := range menuTree {
		// 如果菜单有权限，添加菜单ID
		if menu.HasPermission {
			*menuIDs = append(*menuIDs, menu.ID)

			// 检查菜单的权限列表
			for _, auth := range menu.Meta.AuthList {
				if auth.HasPermission {
					*authIDs = append(*authIDs, auth.ID)
				}
			}
		}

		// 递归处理子菜单
		if len(menu.Children) > 0 {
			extractPermissions(menu.Children, menuIDs, authIDs)
		}
	}
}
