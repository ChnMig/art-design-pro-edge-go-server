package user

import (
	"errors"

	"api-server/api/response"
	userdomain "api-server/domain/admin/user"
	"api-server/util/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReturnDomainError 将 domain 层错误映射为统一的接口错误响应。
func ReturnDomainError(c *gin.Context, err error, fallback string) {
	log.WithRequest(c).Error("用户领域错误", zap.Error(err))

	switch {
	case errors.Is(err, userdomain.ErrUserNotFound):
		response.ReturnError(c, response.DATA_LOSS, "用户不存在")
	case errors.Is(err, userdomain.ErrRoleNotInTenant):
		response.ReturnError(c, response.PERMISSION_DENIED, "角色不存在或不属于当前租户")
	case errors.Is(err, userdomain.ErrCannotDeleteSuperAdmin):
		response.ReturnError(c, response.DATA_LOSS, "不能删除超级管理员")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}

