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
	case errors.Is(err, roledomain.ErrRoleInUse):
		response.ReturnError(c, response.FAILED_PRECONDITION, "角色仍被用户使用，不能禁用或删除")
	case errors.Is(err, roledomain.ErrCannotDeleteTenantAdmin):
		response.ReturnError(c, response.FAILED_PRECONDITION, "租户管理员角色不能删除")
	case errors.Is(err, roledomain.ErrRoleTenantImmutable):
		response.ReturnError(c, response.INVALID_ARGUMENT, "角色归属租户不能修改")
	case errors.Is(err, roledomain.ErrInvalidRoleInput):
		response.ReturnError(c, response.INVALID_ARGUMENT, "角色参数不符合约束")
	case errors.Is(err, roledomain.ErrRoleTenantNotFound):
		response.ReturnError(c, response.NOT_FOUND, "角色所属租户不存在")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}
