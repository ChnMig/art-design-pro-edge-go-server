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
	log.WithRequest(c).Error("平台菜单领域错误", zap.Error(err))

	switch {
	case errors.Is(err, menudomain.ErrParentMenuNotFound):
		response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
	case errors.Is(err, menudomain.ErrParentMenuDisabled):
		response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
	case errors.Is(err, menudomain.ErrDisableMenuWithEnabledChild):
		response.ReturnError(c, response.DATA_LOSS, "请先禁用子菜单")
	case errors.Is(err, menudomain.ErrMenuNotFound):
		response.ReturnError(c, response.DATA_LOSS, "菜单不存在")
	case errors.Is(err, menudomain.ErrMenuHasChildren):
		response.ReturnError(c, response.DATA_LOSS, "请先删除子菜单")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}

