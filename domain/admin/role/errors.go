package role

import "errors"

var (
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")
	// ErrRoleInUse 角色仍被用户使用，不能禁用或删除
	ErrRoleInUse = errors.New("role in use")
	// ErrCannotDeleteTenantAdmin 管理员角色不能删除
	ErrCannotDeleteTenantAdmin = errors.New("cannot delete tenant admin role")
	// ErrRoleTenantImmutable 已创建角色不能跨租户移动
	ErrRoleTenantImmutable = errors.New("role tenant is immutable")
	// ErrInvalidRoleInput 角色参数不符合约束
	ErrInvalidRoleInput = errors.New("invalid role input")
	// ErrRoleTenantNotFound 角色所属租户不存在
	ErrRoleTenantNotFound = errors.New("role tenant not found")
)
