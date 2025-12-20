package health

// StatusDTO 对外返回的健康检查 DTO（Data Transfer Object，数据传输对象）
type StatusDTO struct {
	Status    string `json:"status"`
	Ready     bool   `json:"ready"`
	Uptime    string `json:"uptime"`
	Timestamp int64  `json:"timestamp"`
}
