package menu

import (
	"sort"

	"api-server/db/pgdb/system"
)

// MenuResponse 定义返回给前端的菜单结构
type MenuResponse struct {
	ID        uint           `json:"id"`
	UpdatedAt uint           `json:"updatedAt,omitempty"`
	Path      string         `json:"path"`
	Name      string         `json:"name"`
	Component string         `json:"component,omitempty"`
	Meta      MenuMeta       `json:"meta"`
	Children  []MenuResponse `json:"children,omitempty"`
	ParentID  uint           `json:"parentId,omitempty"`
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
	IsEnable          bool                 `json:"isEnable,omitempty"`
	Sort              uint                 `json:"sort,omitempty"`
	AuthList          []MenuPermissionResp `json:"authList,omitempty"`
}

// MenuPermissionResp 定义菜单权限响应结构
type MenuPermissionResp struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	AuthMark string `json:"auth_mark"` // 对应 MenuPermission 中的 Mark
}

// 构建菜单树
func BuildMenuTree(menus []system.Menu, permissions []system.MenuPermission, all bool) []MenuResponse {
	var menuTree []MenuResponse
	// 找出所有顶级菜单(ParentID = 0)
	var rootMenus []system.Menu
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

// 将数据库菜单转换为响应结构
func convertMenuToResponse(menu system.Menu) MenuResponse {
	return MenuResponse{
		ID:        menu.ID,
		UpdatedAt: uint(menu.UpdatedAt.Unix()),
		Path:      menu.Path,
		Name:      menu.Name,
		Component: menu.Component,
		ParentID:  menu.ParentID,
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
			IsEnable:          menu.Status == 1,            // 1表示启用
			Sort:              menu.Sort,
		},
	}
}

// 递归构建菜单子项
func buildMenuChildren(parent *MenuResponse, allMenus []system.Menu, allPermissions []system.MenuPermission, all bool) {
	var childMenus []system.Menu
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
				child.Meta.AuthList = append(child.Meta.AuthList, MenuPermissionResp{
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

// 对菜单切片按 Sort 从大到小排序，Sort 为 0 则按 ID 从大到小排序
func sortMenus(menus []system.Menu) {
	sort.Slice(menus, func(i, j int) bool {
		// 如果两个菜单都有 Sort 值并且不为 0
		if menus[i].Sort > 0 && menus[j].Sort > 0 {
			// Sort 值不同时，按 Sort 从大到小排序
			if menus[i].Sort != menus[j].Sort {
				return menus[i].Sort > menus[j].Sort // 从大到小排序
			}
			// Sort 值相同时，按 ID 从大到小排序
			return menus[i].ID > menus[j].ID
		}
		// 如果只有其中一个有 Sort 值
		if menus[i].Sort > 0 && menus[j].Sort == 0 {
			return true // i 排在前面
		}
		if menus[i].Sort == 0 && menus[j].Sort > 0 {
			return false // j 排在前面
		}
		// 都没有 Sort 值，按 ID 从大到小排序
		return menus[i].ID > menus[j].ID
	})
}
