package menu

import (
	"errors"

	"api-server/api/response"
	menudomain "api-server/domain/admin/menu"
	"api-server/util/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReturnDomainError 将 domain 层错误映射为统一的接口错误响应。
func ReturnDomainError(c *gin.Context, err error, fallback string) {
	log.WithRequest(c).Error("菜单领域错误", zap.Error(err))

	switch {
	case errors.Is(err, menudomain.ErrRoleNotFound):
		response.ReturnError(c, response.DATA_LOSS, "角色不存在")
	case errors.Is(err, menudomain.ErrPermissionDenied):
		response.ReturnError(c, response.PERMISSION_DENIED, "无权操作该角色菜单")
	case errors.Is(err, menudomain.ErrMenuOutOfScope):
		response.ReturnError(c, response.PERMISSION_DENIED, "菜单超出可分配范围")
	case errors.Is(err, menudomain.ErrAuthOutOfScope):
		response.ReturnError(c, response.PERMISSION_DENIED, "按钮权限超出可分配范围")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}

