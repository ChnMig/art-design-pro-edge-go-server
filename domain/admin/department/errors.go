package department

import "errors"

var (
	// ErrDepartmentNotFound 部门不存在
	ErrDepartmentNotFound = errors.New("department not found")
	// ErrDepartmentHasUsers 部门下仍有用户
	ErrDepartmentHasUsers = errors.New("department has users")
)

