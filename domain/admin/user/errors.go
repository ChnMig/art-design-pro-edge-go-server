package user

import "errors"

var (
	// ErrInvalidCredentials 登录账号或密码错误
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserDisabled 用户已被禁用
	ErrUserDisabled = errors.New("user disabled")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrRoleNotInTenant 角色不存在或不属于当前租户
	ErrRoleNotInTenant = errors.New("role not in tenant")
	// ErrCannotDeleteSuperAdmin 不能删除超级管理员
	ErrCannotDeleteSuperAdmin = errors.New("cannot delete super admin")
	// ErrTenantQueryTooShort 登录页租户搜索输入过短
	ErrTenantQueryTooShort = errors.New("tenant query too short")
)
