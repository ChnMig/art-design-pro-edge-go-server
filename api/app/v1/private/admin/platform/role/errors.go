package role

import (
	"errors"

	"api-server/api/response"
	roledomain "api-server/domain/admin/role"
	"api-server/util/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReturnDomainError 将 domain 层错误映射为统一的接口错误响应。
func ReturnDomainError(c *gin.Context, err error, fallback string) {
	log.WithRequest(c).Error("平台角色领域错误", zap.Error(err))

	switch {
	case errors.Is(err, roledomain.ErrRoleNotFound):
		response.ReturnError(c, response.DATA_LOSS, "角色不存在")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}

