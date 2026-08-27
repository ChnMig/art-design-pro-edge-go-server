package tenant

import (
	"errors"

	"api-server/api/response"
	tenantdomain "api-server/domain/admin/tenant"
	"api-server/util/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReturnDomainError 将 domain 层错误映射为统一的接口错误响应。
func ReturnDomainError(c *gin.Context, err error, fallback string) {
	log.WithRequest(c).Error("租户领域错误", zap.Error(err))

	switch {
	case errors.Is(err, tenantdomain.ErrTenantNotFound):
		response.ReturnError(c, response.DATA_LOSS, "租户不存在")
	case errors.Is(err, tenantdomain.ErrDefaultTenantProtected):
		response.ReturnError(c, response.FAILED_PRECONDITION, "默认平台租户编码不能修改，且不能删除或禁用")
	case errors.Is(err, tenantdomain.ErrInvalidTenantInput):
		response.ReturnError(c, response.INVALID_ARGUMENT, "租户参数不符合安全要求")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}
