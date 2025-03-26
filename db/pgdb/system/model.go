package system

import (
	"gorm.io/gorm"
)

// Department 部门表
type Department struct {
	gorm.Model
	Name   string `json:"name,omitempty"`
	Sort   uint   `json:"sort,omitempty"`
	Status uint   `json:"status,omitempty"` // 状态(1:启用 2:禁用)
	Users  []User `json:"users,omitempty"`  // 一对多关联用户表
}

// Role 角色表
type Role struct {
	gorm.Model
	Name     string     `json:"name,omitempty"`
	Desc     string     `json:"desc,omitempty"`
	Status   uint       `json:"status,omitempty"`                                 // 状态(1:启用 2:禁用)
	Menus    []Menu     `json:"menus,omitempty" gorm:"many2many:role_menu;"`      // 多对多关联菜单表
	User     []User     `json:"users,omitempty"`                                  // 一对多关联用户表
	MenuAuth []MenuAuth `json:"menu_auths,omitempty" gorm:"many2many:role_auth;"` // 多对多关联菜单按钮权限表
}

// Menu 菜单表
type Menu struct {
	gorm.Model
	Path              string     `json:"path,omitempty"`
	Name              string     `json:"name,omitempty"`
	Component         string     `json:"component,omitempty"`            // vue组件
	Title             string     `json:"title,omitempty"`                // 菜单标题
	Icon              string     `json:"icon,omitempty"`                 // 菜单图标
	ShowBadge         uint       `json:"show_badge,omitempty"`           // 是否显示角标(1:显示 2:隐藏)
	ShowTextBadge     string     `json:"show_text_badge,omitempty"`      // 是否显示文本角标(1:显示 2:隐藏)
	IsHide            uint       `json:"is_hide,omitempty"`              // 是否隐藏(1:隐藏 2:显示)
	IsHideTab         uint       `json:"is_hide_tab,omitempty"`          // 是否隐藏标签(1:隐藏 2:显示)
	Link              string     `json:"link,omitempty"`                 // 链接(外链)
	IsIframe          uint       `json:"is_iframe,omitempty"`            // 是否内嵌(1:内嵌 2:不内嵌)
	KeepAlive         uint       `json:"keep_alive,omitempty"`           // 是否缓存(1:缓存 2:不缓存)
	IsInMainContainer uint       `json:"is_in_main_container,omitempty"` // 是否在主容器内(一级菜单使用)(1:是 2:否)
	Status            uint       `json:"status,omitempty"`               // 状态(1:启用 2:禁用)
	Level             uint       `json:"level,omitempty"`                // 层级(从1开始)
	ParentID          uint       `json:"parent_id,omitempty"`            // 父级ID
	Sort              uint       `json:"sort,omitempty"`                 // 排序(从大到小)
	Roles             []Role     `json:"roles,omitempty" gorm:"many2many:role_menu;"`
	MenuAuth          []MenuAuth `json:"menu_auths,omitempty"` // 一对多关联菜单按钮权限表
}

// MenuPermission 菜单按钮权限表
type MenuAuth struct {
	gorm.Model
	MenuID uint   `json:"menu_id,omitempty"`
	Mark   string `json:"mark,omitempty"` // 标识
	Title  string `json:"title,omitempty"`
	Roles  []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"` // 多对多关联角色表
}

// User 用户表
type User struct {
	gorm.Model
	DepartmentID uint   `json:"department_id,omitempty"`
	RoleID       uint   `json:"role_id,omitempty"`
	Name         string `json:"name,omitempty"` // 姓名, 不可修改
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Gender       uint   `json:"gender,omitempty"`
	Status       uint   `json:"status,omitempty"` // 状态(1:启用 2:禁用)
}
