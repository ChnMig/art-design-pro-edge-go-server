package system

import (
	"gorm.io/gorm"
)

type Info struct {
	gorm.Model
}

// Department 部门表
type Department struct {
	gorm.Model
	Name   string `json:"name"`
	Sort   uint   `json:"sort"`
	Status uint   `json:"status"` // 状态(1:启用 2:禁用)
	Users  []User `json:"users"`  // 一对多关联用户表
}

// Role 角色表
type Role struct {
	gorm.Model
	Name            string           `json:"name"`
	Desc            string           `json:"desc"`
	Status          uint             `json:"status"`                                              // 状态(1:启用 2:禁用)
	Menus           []Menu           `json:"menus" gorm:"many2many:role_menu;"`                   // 多对多关联菜单表
	User            []User           `json:"users"`                                               // 一对多关联用户表
	MenuPermissions []MenuPermission `json:"menu_permissions" gorm:"many2many:role_permissions;"` // 一对多关联菜单按钮权限表
}

// Menu 菜单表
type Menu struct {
	gorm.Model
	Path              string           `json:"path"`
	Name              string           `json:"name"`
	Component         string           `json:"component"`            // vue组件
	Title             string           `json:"title"`                // 菜单标题
	Icon              string           `json:"icon"`                 // 菜单图标
	ShowBadge         uint             `json:"show_badge"`           // 是否显示角标(1:显示 2:隐藏)
	ShowTextBadge     string           `json:"show_text_badge"`      // 是否显示文本角标(1:显示 2:隐藏)
	IsHide            uint             `json:"is_hide"`              // 是否隐藏(1:隐藏 2:显示)
	IsHideTab         uint             `json:"is_hide_tab"`          // 是否隐藏标签(1:隐藏 2:显示)
	Link              string           `json:"link"`                 // 链接(外链)
	IsIframe          uint             `json:"is_iframe"`            // 是否内嵌(1:内嵌 2:不内嵌)
	KeepAlive         uint             `json:"keep_alive"`           // 是否缓存(1:缓存 2:不缓存)
	IsInMainContainer uint             `json:"is_in_main_container"` // 是否在主容器内(1:是 2:否)
	Status            uint             `json:"status"`               // 状态(1:启用 2:禁用)
	Level             uint             `json:"level"`                // 层级(从1开始)
	ParentID          uint             `json:"parent_id"`            // 父级ID
	Roles             []Role           `json:"roles" gorm:"many2many:role_menu;"`
	MenuPermissions   []MenuPermission `json:"menu_permissions"` // 一对多关联菜单按钮权限表
}

// MenuPermission 菜单按钮权限表
type MenuPermission struct {
	gorm.Model
	MenuID uint   `json:"menu_id"`
	Mark   string `json:"mark"` // 标识
	Title  string `json:"title"`
	Sort   uint   `json:"sort"`
	Roles  []Role `json:"roles" gorm:"many2many:role_permissions;"` // 多对多关联角色表
}

// User 用户表
type User struct {
	gorm.Model
	DepartmentID uint   `json:"department_id"`
	RoleID       uint   `json:"role_id"`
	Name         string `json:"name"` // 姓名, 不可修改
	Username     string `json:"username"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	Gender       uint   `json:"gender"`
	Status       uint   `json:"status"` // 状态(1:启用 2:禁用)
}
