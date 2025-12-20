package health

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-server/api/response"
	domain "api-server/domain/health"
	"api-server/util/log"
)

// Status 返回服务健康状态
func Status(c *gin.Context) {
	// 默认使用带 trace_id 等上下文信息的 logger
	l := log.FromContext(c)
	l.Debug("健康检查开始")

	status, err := domain.GetStatus()
	if err != nil {
		log.WithRequest(c).Error("健康检查失败", zap.Error(err))
		ReturnDomainError(c, err)
		return
	}

	dto := StatusDTO{
		Status:    status.Status,
		Ready:     status.Ready,
		Uptime:    status.Uptime.String(),
		Timestamp: status.Timestamp,
	}
	response.ReturnData(c, dto)
}
