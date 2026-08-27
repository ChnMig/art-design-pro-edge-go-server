package tenant

import "errors"

var (
	// ErrTenantNotFound 租户不存在
	ErrTenantNotFound = errors.New("tenant not found")
	// ErrDefaultTenantProtected 默认平台租户不能被删除或禁用
	ErrDefaultTenantProtected = errors.New("default tenant is protected")
	// ErrInvalidTenantInput 租户参数不符合安全要求
	ErrInvalidTenantInput = errors.New("invalid tenant input")
)
