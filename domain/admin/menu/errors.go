package menu

import "errors"

var (
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")
	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = errors.New("permission denied")

	// ErrMenuNotFound 菜单不存在
	ErrMenuNotFound = errors.New("menu not found")
	// ErrParentMenuNotFound 父级菜单不存在
	ErrParentMenuNotFound = errors.New("parent menu not found")
	// ErrParentMenuDisabled 父级菜单已禁用
	ErrParentMenuDisabled = errors.New("parent menu disabled")
	// ErrMenuHasChildren 菜单存在子菜单
	ErrMenuHasChildren = errors.New("menu has children")
	// ErrDisableMenuWithEnabledChild 请先禁用子菜单
	ErrDisableMenuWithEnabledChild = errors.New("disable menu with enabled child")

	// ErrMenuOutOfScope 菜单超出可分配范围
	ErrMenuOutOfScope = errors.New("menu out of scope")
	// ErrAuthOutOfScope 按钮权限超出可分配范围
	ErrAuthOutOfScope = errors.New("auth out of scope")
)

