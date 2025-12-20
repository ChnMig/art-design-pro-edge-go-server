package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-server/api/auth"
	"api-server/api/middleware"
	"api-server/api/response"
	userdomain "api-server/domain/admin/user"
)

func Login(c *gin.Context) {
	params := &struct {
		TenantCode string `json:"tenant_code" form:"tenant_code" binding:"required"` // 企业编号
		Account    string `json:"account" form:"account" binding:"required"`         // 登录账号
		Password   string `json:"password" form:"password" binding:"required"`
		Captcha    string `json:"captcha" form:"captcha" binding:"required"`
		CaptchaID  string `json:"captcha_id" form:"captcha_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	// 验证验证码
	if !userdomain.VerifyCaptcha(params.CaptchaID, params.Captcha) {
		response.ReturnError(c, response.INVALID_ARGUMENT, "验证码错误")
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 查询用户（多租户验证）
	user, tenant, err := userdomain.VerifyLogin(userdomain.LoginInput{
		TenantCode: params.TenantCode,
		Account:    params.Account,
		Password:   params.Password,
	})

	if err != nil {
		_ = userdomain.CreateLoginLogFromInput(userdomain.LoginLogInput{
			TenantCode:  params.TenantCode,
			UserName:    params.Account,
			IP:          clientIP,
			LoginStatus: "failed",
		})

		switch {
		case errors.Is(err, userdomain.ErrInvalidCredentials):
			response.ReturnError(c, response.INVALID_ARGUMENT, "账号或密码错误")
		case errors.Is(err, userdomain.ErrUserDisabled):
			response.ReturnError(c, response.INVALID_ARGUMENT, "账号已被禁用")
		default:
			zap.L().Error("查询用户失败", zap.Error(err))
			response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		}
		return
	}

	// 记录登录成功日志
	_ = userdomain.CreateLoginLogFromInput(userdomain.LoginLogInput{
		TenantCode:  params.TenantCode,
		UserName:    params.Account,
		IP:          clientIP,
		LoginStatus: "success",
	})
	// 生成多租户token
	token, err := auth.JWTIssue(user.ID, tenant.ID, user.Account)
	if err != nil {
		zap.L().Error("生成token失败", zap.Error(err))
		response.ReturnError(c, response.INTERNAL, "生成token失败")
		return
	}
	response.ReturnData(c, gin.H{
		"access_token": token,
		"tenant_info": gin.H{
			"tenant_id":   tenant.ID,
			"tenant_code": tenant.Code,
			"tenant_name": tenant.Name,
		},
		"user_info": gin.H{
			"user_id":  user.ID,
			"name":     user.Name,
			"username": user.Username,
			"account":  user.Account,
		},
	})
}

func FindLoginLogList(c *gin.Context) {
	params := &struct {
		IP       string `json:"ip" form:"ip"`
		Username string `json:"username" form:"username"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	// 获取分页参数
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)
	logs, total, err := userdomain.FindLoginLogList(userdomain.FindLoginLogQuery{
		IP:       params.IP,
		Username: params.Username,
	}, page, pageSize)
	if err != nil {
		zap.L().Error("查询登录日志失败", zap.Error(err))
		response.ReturnError(c, response.DATA_LOSS, "查询登录日志失败")
		return
	}
	response.ReturnDataWithTotal(c, int(total), logs)
}
