package health

import "errors"

// 领域层健康检查相关错误定义
// 仅描述业务语义，不关心具体返回给前端的文案
var (
	// ErrServiceNotReady 表示服务尚未就绪，暂不可对外提供服务
	ErrServiceNotReady = errors.New("service not ready")

	// ErrServiceUnhealthy 表示服务健康检查失败
	ErrServiceUnhealthy = errors.New("service unhealthy")
)
