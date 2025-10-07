package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-server/api/auth"
	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
	"api-server/db/rdb/captcha"
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
	captchaVerify := captcha.GetRedisStore().Verify(params.CaptchaID, params.Captcha, true)
	if !captchaVerify {
		response.ReturnError(c, response.INVALID_ARGUMENT, "验证码错误")
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 查询用户（多租户验证）
	user, tenant, err := system.VerifyUser(params.TenantCode, params.Account, params.Password)
	if err != nil {
		zap.L().Error("查询用户失败", zap.Error(err))
		// 记录登录失败日志（验证码正确但查询失败）
		failedLog := system.SystemUserLoginLog{
			TenantCode:  params.TenantCode,
			UserName:    params.Account,
			Password:    "", // 安全：不记录密码
			IP:          clientIP,
			LoginStatus: "failed",
		}
		system.CreateLoginLog(&failedLog)
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}
	if user.ID == 0 {
		// 记录登录失败日志（账号或密码错误）
		failedLog := system.SystemUserLoginLog{
			TenantCode:  params.TenantCode,
			UserName:    params.Account,
			Password:    "", // 安全：不记录密码
			IP:          clientIP,
			LoginStatus: "failed",
		}
		system.CreateLoginLog(&failedLog)
		response.ReturnError(c, response.INVALID_ARGUMENT, "账号或密码错误")
		return
	}
	if user.Status != system.StatusEnabled {
		// 记录登录失败日志（账号已被禁用）
		failedLog := system.SystemUserLoginLog{
			TenantCode:  params.TenantCode,
			UserName:    params.Account,
			Password:    "", // 安全：不记录密码
			IP:          clientIP,
			LoginStatus: "failed",
		}
		system.CreateLoginLog(&failedLog)
		response.ReturnError(c, response.INVALID_ARGUMENT, "账号已被禁用")
		return
	}

	// 记录登录成功日志
	successLog := system.SystemUserLoginLog{
		TenantCode:  params.TenantCode,
		UserName:    params.Account,
		Password:    "", // 安全：不记录密码
		IP:          clientIP,
		LoginStatus: "success",
	}
	system.CreateLoginLog(&successLog)
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
	log := system.SystemUserLoginLog{
		IP:       params.IP,
		UserName: params.Username,
	}
	logs, total, err := system.FindLoginLogList(&log, page, pageSize)
	if err != nil {
		zap.L().Error("查询登录日志失败", zap.Error(err))
		response.ReturnError(c, response.DATA_LOSS, "查询登录日志失败")
		return
	}
	response.ReturnDataWithTotal(c, int(total), logs)
}
