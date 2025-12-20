package menu

import (
	"errors"

	commonmenu "api-server/common/menu"
	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

func GetPlatformMenuTree() ([]commonmenu.MenuResponse, error) {
	menus, allAuths, err := system.GetMenuData()
	if err != nil {
		return nil, err
	}
	return commonmenu.BuildMenuTree(menus, allAuths, true), nil
}

type AddMenuInput struct {
	Path          string
	Name          string
	Component     string
	Title         string
	Icon          string
	ShowBadge     uint
	ShowTextBadge string
	IsHide        uint
	IsHideTab     uint
	Link          string
	IsIframe      uint
	KeepAlive     uint
	IsFirstLevel  uint
	Status        uint
	ParentID      uint
	Sort          uint
}

func AddMenu(input AddMenuInput) (system.SystemMenu, error) {
	if input.ShowBadge == 0 {
		input.ShowBadge = 2
	}

	var level uint = 1
	if input.ParentID != 0 {
		parentMenu := system.SystemMenu{Model: gorm.Model{ID: input.ParentID}}
		if err := system.GetMenu(&parentMenu); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return system.SystemMenu{}, ErrParentMenuNotFound
			}
			return system.SystemMenu{}, err
		}
		if parentMenu.Status != system.StatusEnabled {
			return system.SystemMenu{}, ErrParentMenuDisabled
		}
		level = parentMenu.Level + 1
	}

	menu := system.SystemMenu{
		Path:          input.Path,
		Name:          input.Name,
		Component:     input.Component,
		Title:         input.Title,
		Icon:          input.Icon,
		ShowBadge:     input.ShowBadge,
		ShowTextBadge: input.ShowTextBadge,
		IsHide:        input.IsHide,
		IsHideTab:     input.IsHideTab,
		Link:          input.Link,
		IsIframe:      input.IsIframe,
		KeepAlive:     input.KeepAlive,
		IsFirstLevel:  input.IsFirstLevel,
		Status:        input.Status,
		Level:         level,
		ParentID:      input.ParentID,
		Sort:          input.Sort,
	}
	if err := system.AddMenu(&menu); err != nil {
		return system.SystemMenu{}, err
	}
	return menu, nil
}

type UpdateMenuInput struct {
	ID            uint
	Path          string
	Name          string
	Component     string
	Title         string
	Icon          string
	ShowBadge     uint
	ShowTextBadge string
	IsHide        uint
	IsHideTab     uint
	Link          string
	IsIframe      uint
	KeepAlive     uint
	IsFirstLevel  uint
	Status        uint
	ParentID      uint
	Sort          uint
}

func UpdateMenu(input UpdateMenuInput) (system.SystemMenu, error) {
	if input.ShowBadge == 0 {
		input.ShowBadge = 2
	}

	var level uint = 1
	if input.ParentID != 0 {
		parent := system.SystemMenu{Model: gorm.Model{ID: input.ParentID}}
		if err := system.GetMenu(&parent); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return system.SystemMenu{}, ErrParentMenuNotFound
			}
			return system.SystemMenu{}, err
		}
		if parent.Status != system.StatusEnabled {
			return system.SystemMenu{}, ErrParentMenuDisabled
		}
		level = parent.Level + 1
	}

	if input.Status == system.StatusDisabled {
		children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: input.ID}, -1, -1)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemMenu{}, err
		}
		for _, child := range children {
			if child.Status == system.StatusEnabled {
				return system.SystemMenu{}, ErrDisableMenuWithEnabledChild
			}
		}
	}

	menu := system.SystemMenu{
		Model:         gorm.Model{ID: input.ID},
		Path:          input.Path,
		Name:          input.Name,
		Component:     input.Component,
		Title:         input.Title,
		Icon:          input.Icon,
		ShowBadge:     input.ShowBadge,
		ShowTextBadge: input.ShowTextBadge,
		IsHide:        input.IsHide,
		IsHideTab:     input.IsHideTab,
		Link:          input.Link,
		IsIframe:      input.IsIframe,
		KeepAlive:     input.KeepAlive,
		IsFirstLevel:  input.IsFirstLevel,
		Status:        input.Status,
		Level:         level,
		ParentID:      input.ParentID,
		Sort:          input.Sort,
	}
	if err := system.UpdateMenu(&menu); err != nil {
		return system.SystemMenu{}, err
	}
	return menu, nil
}

func DeleteMenu(id uint) (system.SystemMenu, error) {
	menu := system.SystemMenu{Model: gorm.Model{ID: id}}
	if err := system.GetMenu(&menu); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemMenu{}, ErrMenuNotFound
		}
		return system.SystemMenu{}, err
	}

	children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: menu.ID}, -1, -1)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return system.SystemMenu{}, err
	}
	if len(children) > 0 {
		return system.SystemMenu{}, ErrMenuHasChildren
	}

	if err := system.DeleteMenu(&menu); err != nil {
		return system.SystemMenu{}, err
	}
	return menu, nil
}

func GetMenuAuthList(menuID uint) ([]system.SystemMenuAuth, error) {
	auth := system.SystemMenuAuth{MenuID: menuID}
	return system.FindMenuAuthList(&auth)
}

type AddMenuAuthInput struct {
	MenuID uint
	Mark   string
	Title  string
}

func AddMenuAuth(input AddMenuAuthInput) (system.SystemMenuAuth, error) {
	auth := system.SystemMenuAuth{
		MenuID: input.MenuID,
		Mark:   input.Mark,
		Title:  input.Title,
	}
	if err := system.AddMenuAuth(&auth); err != nil {
		return system.SystemMenuAuth{}, err
	}
	return auth, nil
}

type UpdateMenuAuthInput struct {
	ID     uint
	Title  string
	Mark   string
	MenuID uint
}

func UpdateMenuAuth(input UpdateMenuAuthInput) (system.SystemMenuAuth, error) {
	auth := system.SystemMenuAuth{
		Model:  gorm.Model{ID: input.ID},
		Title:  input.Title,
		Mark:   input.Mark,
		MenuID: input.MenuID,
	}
	if err := system.UpdateMenuAuth(&auth); err != nil {
		return system.SystemMenuAuth{}, err
	}
	return auth, nil
}

func DeleteMenuAuth(id uint) (system.SystemMenuAuth, error) {
	auth := system.SystemMenuAuth{Model: gorm.Model{ID: id}}
	if err := system.DeleteMenuAuth(&auth); err != nil {
		return system.SystemMenuAuth{}, err
	}
	return auth, nil
}

func GetTenantMenuTree(tenantID uint) ([]commonmenu.MenuResponse, error) {
	menus, allAuths, err := system.GetMenuData()
	if err != nil {
		return nil, err
	}
	scopeIDs, err := system.GetTenantMenuScopeIDs(tenantID)
	if err != nil {
		return nil, err
	}
	authScopeIDs, err := system.GetTenantAuthScopeIDs(tenantID)
	if err != nil {
		return nil, err
	}
	return commonmenu.BuildMenuTreeWithPermission(menus, allAuths, scopeIDs, authScopeIDs, true), nil
}

func UpdateTenantMenuScope(tenantID uint, menuData []commonmenu.MenuResponse) ([]commonmenu.MenuResponse, error) {
	menuIDs := extractCheckedMenuIDs(menuData)
	authIDs := extractCheckedAuthIDs(menuData)

	if err := system.SaveTenantMenuScope(tenantID, menuIDs); err != nil {
		return nil, err
	}
	if err := system.SaveTenantAuthScope(tenantID, authIDs); err != nil {
		return nil, err
	}
	if err := system.PruneTenantRoleAssociations(tenantID, menuIDs, authIDs); err != nil {
		return nil, err
	}

	menus, allAuths, err := system.GetMenuData()
	if err != nil {
		return nil, err
	}
	return commonmenu.BuildMenuTreeWithPermission(menus, allAuths, menuIDs, authIDs, true), nil
}

