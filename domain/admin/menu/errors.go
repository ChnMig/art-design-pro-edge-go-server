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
	// ErrMenuCycle 菜单父子关系形成环
	ErrMenuCycle = errors.New("menu hierarchy cycle")
	// ErrMenuAuthNotFound 菜单按钮权限不存在
	ErrMenuAuthNotFound = errors.New("menu auth not found")
	// ErrMenuAuthMenuImmutable 按钮权限创建后不能迁移到其他菜单
	ErrMenuAuthMenuImmutable = errors.New("menu auth menu is immutable")
	// ErrInvalidMenuStatus 菜单状态只能为启用或禁用
	ErrInvalidMenuStatus = errors.New("invalid menu status")

	// ErrMenuOutOfScope 菜单超出可分配范围
	ErrMenuOutOfScope = errors.New("menu out of scope")
	// ErrAuthOutOfScope 按钮权限超出可分配范围
	ErrAuthOutOfScope = errors.New("auth out of scope")
	// ErrAuthWithoutMenu 按钮权限不能脱离所属菜单单独授权
	ErrAuthWithoutMenu = errors.New("auth requires selected menu")
)
