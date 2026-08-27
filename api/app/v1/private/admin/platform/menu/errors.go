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
	case errors.Is(err, menudomain.ErrMenuCycle):
		response.ReturnError(c, response.INVALID_ARGUMENT, "菜单不能将自身或后代设为父级")
	case errors.Is(err, menudomain.ErrMenuAuthNotFound):
		response.ReturnError(c, response.NOT_FOUND, "菜单按钮权限不存在")
	case errors.Is(err, menudomain.ErrMenuAuthMenuImmutable):
		response.ReturnError(c, response.INVALID_ARGUMENT, "按钮权限创建后不能迁移到其他菜单")
	case errors.Is(err, menudomain.ErrInvalidMenuStatus):
		response.ReturnError(c, response.INVALID_ARGUMENT, "菜单状态只能为启用或禁用")
	case errors.Is(err, menudomain.ErrAuthWithoutMenu):
		response.ReturnError(c, response.INVALID_ARGUMENT, "按钮权限必须与所属菜单一起授权")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}
