package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
	"api-server/db/rdb/captcha"
	"api-server/util/authentication"
)

func Login(c *gin.Context) {
	params := &struct {
		Username  string `json:"username" form:"username" binding:"required"`
		Password  string `json:"password" form:"password" binding:"required"`
		Captcha   string `json:"captcha" form:"captcha" binding:"required"`
		CaptchaID string `json:"captcha_id" form:"captcha_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	// 验证验证码
	captchaVerify := captcha.GetRedisStore().Verify(params.CaptchaID, params.Captcha, true)
	if !captchaVerify {
		response.ReturnError(c, response.INVALID_ARGUMENT, "验证码错误")
		return
	}
	// 查询用户
	user, err := system.VerifyUser(params.Username, params.Password)
	if err != nil {
		zap.L().Error("查询用户失败", zap.Error(err))
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}
	if user.ID == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "账号或密码错误")
		return
	}
	if user.Status != 1 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "账号已被禁用")
		return
	}
	// 生成token
	token, err := authentication.JWTIssue(fmt.Sprintf("%d", user.ID))
	if err != nil {
		zap.L().Error("生成token失败", zap.Error(err))
		response.ReturnError(c, response.INTERNAL, "生成token失败")
		return
	}
	response.ReturnOk(c, gin.H{
		"access_token": token,
	})
}
