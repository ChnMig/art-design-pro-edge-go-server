package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-server/api/middleware"
	"api-server/api/response"
	userdomain "api-server/domain/admin/user"
)

// SearchTenantCodeForLogin 登录页用：根据输入模糊查询租户编码，返回最多10条
func SearchTenantCodeForLogin(c *gin.Context) {
	params := &struct {
		Code string `json:"code" form:"code" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	result, err := userdomain.SuggestTenantForLogin(params.Code, 10)
	if err != nil {
		if errors.Is(err, userdomain.ErrTenantQueryTooShort) {
			response.ReturnError(c, response.INVALID_ARGUMENT, "输入长度过短")
			return
		}
		zap.L().Error("查询租户编码建议失败", zap.Error(err))
		response.ReturnError(c, response.DATA_LOSS, "查询失败")
		return
	}

	response.ReturnData(c, result)
}
