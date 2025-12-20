package health

import "time"

// Status 领域层的健康状态实体
// 不关心具体展示格式，仅承载核心健康信息
type Status struct {
	Status    string
	Ready     bool
	Uptime    time.Duration
	Timestamp int64
}

var startTime = time.Now()

// GetStatus 获取当前服务的健康状态（领域层示例）
// 这里简单基于进程启动时间和当前时间构造状态，
// 当前示例中始终返回 nil 错误，便于在 handler 中演示错误映射用法。
func GetStatus() (Status, error) {
	return Status{
		Status:    "ok",
		Ready:     true,
		Uptime:    time.Since(startTime),
		Timestamp: time.Now().Unix(),
	}, nil
}
