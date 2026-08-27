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
	// ErrDepartmentNotInTenant 部门不存在或不属于当前租户
	ErrDepartmentNotInTenant = errors.New("department not in tenant")
	// ErrUserNotInTenant 用户不存在或不属于当前租户
	ErrUserNotInTenant = errors.New("user not in tenant")
	// ErrInvalidTenantContext 当前请求缺少有效租户上下文
	ErrInvalidTenantContext = errors.New("invalid tenant context")
	// ErrUnsafeCacheData 缓存缺少租户归属信息，不能安全返回
	ErrUnsafeCacheData = errors.New("unsafe cache data")
	// ErrCannotDeleteSuperAdmin 不能删除超级管理员
	ErrCannotDeleteSuperAdmin = errors.New("cannot delete super admin")
	// ErrCannotModifySuperAdmin 非平台超级管理员不能修改平台根用户
	ErrCannotModifySuperAdmin = errors.New("cannot modify super admin")
	// ErrTenantQueryTooShort 登录页租户搜索输入过短
	ErrTenantQueryTooShort = errors.New("tenant query too short")
	// ErrInvalidCaptchaSize 验证码图片尺寸超出允许范围
	ErrInvalidCaptchaSize = errors.New("invalid captcha size")
	// ErrInvalidUserInput 用户状态或密码不符合安全要求
	ErrInvalidUserInput = errors.New("invalid user input")
)
