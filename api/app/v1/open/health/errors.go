package health

import (
	"errors"

	"api-server/api/response"
	domain "api-server/domain/health"
	"api-server/util/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReturnDomainError 将领域层健康检查错误映射为统一的接口错误响应
func ReturnDomainError(c *gin.Context, err error) {
	log.WithRequest(c).Error("健康检查领域错误", zap.Error(err))

	switch {
	case errors.Is(err, domain.ErrServiceNotReady):
		response.ReturnError(c, response.FAILED_PRECONDITION, "服务尚未就绪，请稍后重试")
	case errors.Is(err, domain.ErrServiceUnhealthy):
		response.ReturnError(c, response.UNAVAILABLE, "服务当前不可用，请稍后重试")
	default:
		response.ReturnError(c, response.INTERNAL, "服务内部错误")
	}
}
