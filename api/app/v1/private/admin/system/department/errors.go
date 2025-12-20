package department

import (
	"errors"

	"api-server/api/response"
	departmentdomain "api-server/domain/admin/department"
	"api-server/util/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReturnDomainError 将 domain 层错误映射为统一的接口错误响应。
func ReturnDomainError(c *gin.Context, err error, fallback string) {
	log.WithRequest(c).Error("部门领域错误", zap.Error(err))

	switch {
	case errors.Is(err, departmentdomain.ErrDepartmentNotFound):
		response.ReturnError(c, response.DATA_LOSS, "部门不存在")
	case errors.Is(err, departmentdomain.ErrDepartmentHasUsers):
		response.ReturnError(c, response.DATA_LOSS, "请先删除部门下的用户")
	default:
		response.ReturnError(c, response.DATA_LOSS, fallback)
	}
}

