package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func GetMenuList(c *gin.Context) {
	uID := c.GetString(middleware.JWTDataKey)
	if uID == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	id, err := strconv.ParseUint(uID, 10, 64)
	if err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "无效的用户ID")
		return
	}

	// 查询用户菜单数据
	_, _, menus, permissions, err := system.GetUserMenuData(uint(id))
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户菜单失败")
		return
	}

	// 构建菜单树
	menuTree := buildMenuTree(menus, permissions)

	response.ReturnOk(c, menuTree)
}

// 构建菜单树
func buildMenuTree(menus []system.Menu, permissions []system.MenuPermission) []system.MenuResponse {
	var menuTree []system.MenuResponse

	// 找出所有顶级菜单(ParentID = 0)
	var rootMenus []system.Menu
	if len(menus) > 0 {
		for _, menu := range menus {
			if menu.ParentID == 0 && menu.Status == 1 {
				rootMenus = append(rootMenus, menu)
			}
		}
	}

	// 递归构建菜单树
	for _, rootMenu := range rootMenus {
		menuResp := convertMenuToResponse(rootMenu)
		buildMenuChildren(&menuResp, menus, permissions)
		menuTree = append(menuTree, menuResp)
	}

	return menuTree
}

// 将数据库菜单转换为响应结构
func convertMenuToResponse(menu system.Menu) system.MenuResponse {
	return system.MenuResponse{
		ID:        menu.ID,
		Path:      menu.Path,
		Name:      menu.Name,
		Component: menu.Component,
		Meta: system.MenuMeta{
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

// 递归构建菜单子项
func buildMenuChildren(parent *system.MenuResponse, allMenus []system.Menu, allPermissions []system.MenuPermission) {
	for _, menu := range allMenus {
		if menu.ParentID == parent.ID && menu.Status == 1 {
			child := convertMenuToResponse(menu)
			// 为子菜单添加权限列表
			for _, perm := range allPermissions {
				if perm.MenuID == menu.ID {
					child.Meta.AuthList = append(child.Meta.AuthList, system.MenuPermissionResp{
						ID:       perm.ID,
						Title:    perm.Title,
						AuthMark: perm.Mark,
					})
				}
			}
			// 递归处理这个子菜单的子菜单
			buildMenuChildren(&child, allMenus, allPermissions)
			parent.Children = append(parent.Children, child)
		}
	}
}
