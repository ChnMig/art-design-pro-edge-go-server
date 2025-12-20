package tenant

import "errors"

var (
	// ErrTenantNotFound 租户不存在
	ErrTenantNotFound = errors.New("tenant not found")
)

