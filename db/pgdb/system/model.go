package system

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户/企业表
type SystemTenant struct {
	gorm.Model
	Code        string             `json:"code,omitempty" gorm:"uniqueIndex;not null"` // 企业编号，唯一
	Name        string             `json:"name,omitempty" gorm:"not null"`             // 企业名称
	Contact     string             `json:"contact,omitempty"`                          // 联系人
	Phone       string             `json:"phone,omitempty"`                            // 联系电话
	Email       string             `json:"email,omitempty"`                            // 邮箱
	Address     string             `json:"address,omitempty"`                          // 地址
	Status      uint               `json:"status,omitempty" gorm:"default:1"`          // 状态(1:启用 2:禁用)
	ExpiredAt   *time.Time         `json:"expired_at,omitempty"`                       // 过期时间
	MaxUsers    uint               `json:"max_users,omitempty" gorm:"default:100"`     // 最大用户数
	SystemUsers []SystemUser       `json:"users,omitempty" gorm:"foreignKey:TenantID"`
	Departments []SystemDepartment `json:"departments,omitempty" gorm:"foreignKey:TenantID"`
	Roles       []SystemRole       `json:"roles,omitempty" gorm:"foreignKey:TenantID"`
}

// Department 部门表
type SystemDepartment struct {
	gorm.Model
	TenantID    uint         `json:"tenant_id,omitempty" gorm:"not null;index"` // 租户ID
	Name        string       `json:"name,omitempty"`
	Sort        uint         `json:"sort,omitempty"`
	Status      uint         `json:"status,omitempty"` // 状态(1:启用 2:禁用)
	SystemUsers []SystemUser `json:"users,omitempty" gorm:"foreignKey:DepartmentID"`
}

// Role 角色表
type SystemRole struct {
	gorm.Model
	TenantID        uint             `json:"tenant_id,omitempty" gorm:"not null;index"` // 租户ID
	Name            string           `json:"name,omitempty"`
	Desc            string           `json:"desc,omitempty"`
	Status          uint             `json:"status,omitempty"`                                             // 状态(1:启用 2:禁用)
	SystemMenus     []SystemMenu     `json:"menus,omitempty" gorm:"many2many:system_roles__system_menus;"` // 多对多关联菜单表
	SystemUsers     []SystemUser     `json:"users,omitempty" gorm:"foreignKey:RoleID"`
	SystemMenuAuths []SystemMenuAuth `json:"menu_auths,omitempty" gorm:"many2many:system_roles__system_auths;"` // 多对多关联菜单按钮权限表
}

// Menu 菜单表
type SystemMenu struct {
	gorm.Model
	Path            string           `json:"path,omitempty"`
	Name            string           `json:"name,omitempty"`
	Component       string           `json:"component,omitempty"`            // vue组件
	Title           string           `json:"title,omitempty"`                // 菜单标题
	Icon            string           `json:"icon,omitempty"`                 // 菜单图标
	ShowBadge       uint             `json:"show_badge,omitempty"`           // 是否显示角标(1:显示 2:隐藏)
	ShowTextBadge   string           `json:"show_text_badge,omitempty"`      // 是否显示文本角标(1:显示 2:隐藏)
	IsHide          uint             `json:"is_hide,omitempty"`              // 是否隐藏(1:隐藏 2:显示)
	IsHideTab       uint             `json:"is_hide_tab,omitempty"`          // 是否隐藏标签(1:隐藏 2:显示)
	Link            string           `json:"link,omitempty"`                 // 链接(外链)
	IsIframe        uint             `json:"is_iframe,omitempty"`            // 是否内嵌(1:内嵌 2:不内嵌)
	KeepAlive       uint             `json:"keep_alive,omitempty"`           // 是否缓存(1:缓存 2:不缓存)
	IsFirstLevel    uint             `json:"is_in_main_container,omitempty"` // 是否在主容器内(一级菜单使用)(1:是 2:否)
	Status          uint             `json:"status,omitempty"`               // 状态(1:启用 2:禁用)
	Level           uint             `json:"level,omitempty"`                // 层级(从1开始)
	ParentID        uint             `json:"parent_id,omitempty"`            // 父级ID
	Sort            uint             `json:"sort,omitempty"`                 // 排序(从小到大，值越小越靠前)
	SystemRoles     []SystemRole     `json:"roles,omitempty" gorm:"many2many:system_roles__system_menus;"`
	SystemMenuAuths []SystemMenuAuth `json:"menu_auths,omitempty" gorm:"foreignKey:MenuID"`
}

// MenuPermission 菜单按钮权限表
type SystemMenuAuth struct {
	gorm.Model
	MenuID      uint         `json:"menu_id,omitempty"`
	Mark        string       `json:"mark,omitempty"` // 标识
	Title       string       `json:"title,omitempty"`
	SystemRoles []SystemRole `json:"roles,omitempty" gorm:"many2many:system_roles__system_auths;"` // 多对多关联角色表
}

// User 用户表
type SystemUser struct {
	gorm.Model
	TenantID     uint   `json:"tenant_id,omitempty" gorm:"not null;index;uniqueIndex:idx_tenant_account"` // 租户ID
	DepartmentID uint   `json:"department_id,omitempty"`
	RoleID       uint   `json:"role_id,omitempty"`
	Name         string `json:"name,omitempty"`                                          // 姓名
	Username     string `json:"username,omitempty"`                                      // 昵称
	Account      string `json:"account,omitempty" gorm:"uniqueIndex:idx_tenant_account"` // 登录账号，同租户内唯一
	Password     string `json:"password,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Gender       uint   `json:"gender,omitempty"` // 性别(1:男 2:女)
	Status       uint   `json:"status,omitempty"` // 状态(1:启用 2:禁用)
}

type SystemUserLoginLog struct {
	gorm.Model
	TenantCode  string `json:"tenant_code,omitempty"` // 企业编号
	UserName    string `json:"user_name,omitempty"`   // 登录账号
	Password    string `json:"password,omitempty"`    // 注意：此字段应为空，不记录实际密码
	IP          string `json:"ip,omitempty"`
	LoginStatus string `json:"login_status,omitempty"` // 登录状态：success, failed
}
