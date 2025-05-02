package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/rdb/captcha"
)

func GetCaptcha(c *gin.Context) {
	params := &struct {
		Width  int `json:"width" form:"width" binding:"required"`
		Height int `json:"height" form:"height" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	driver := base64Captcha.NewDriverDigit(params.Height, params.Width, 6, 0.2, 50)
	client := base64Captcha.NewCaptcha(driver, captcha.GetRedisStore())
	id, b64s, _, err := client.Generate()
	if err != nil {
		response.ReturnError(c, response.UNKNOWN, "验证码生成失败")
		zap.L().Error("验证码生成失败", zap.Error(err))
		return
	}
	response.ReturnData(c, gin.H{
		"id":    id,
		"image": b64s,
	})
}
