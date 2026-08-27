package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-server/api/middleware"
	"api-server/api/response"
	userdomain "api-server/domain/admin/user"
)

func GetCaptcha(c *gin.Context) {
	params := &struct {
		Width  int `json:"width" form:"width" binding:"required"`
		Height int `json:"height" form:"height" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	result, err := userdomain.GenerateCaptcha(params.Width, params.Height)
	if err != nil {
		if errors.Is(err, userdomain.ErrInvalidCaptchaSize) {
			response.ReturnError(c, response.INVALID_ARGUMENT, "验证码宽度需为 80-600，高度需为 40-200")
			return
		}
		response.ReturnError(c, response.UNKNOWN, "验证码生成失败")
		zap.L().Error("验证码生成失败", zap.Error(err))
		return
	}
	response.ReturnData(c, result)
}
