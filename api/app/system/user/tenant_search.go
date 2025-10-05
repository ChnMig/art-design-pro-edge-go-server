package user

import (
    "unicode/utf8"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "api-server/api/middleware"
    "api-server/api/response"
    "api-server/config"
    "api-server/db/pgdb/system"
)

// SearchTenantCodeForLogin 登录页用：根据输入模糊查询租户编码，返回最多10条
func SearchTenantCodeForLogin(c *gin.Context) {
    params := &struct {
        Code string `json:"code" form:"code" binding:"required"`
    }{}
    if !middleware.CheckParam(params, c) {
        return
    }

    // 校验最小长度（必须大于配置值）
    if utf8.RuneCountInString(params.Code) <= config.TenantMinQueryLength {
        response.ReturnError(c, response.INVALID_ARGUMENT, "输入长度过短")
        return
    }

    tenants, err := system.SuggestTenantByCode(params.Code, 10)
    if err != nil {
        zap.L().Error("查询租户编码建议失败", zap.Error(err))
        response.ReturnError(c, response.DATA_LOSS, "查询失败")
        return
    }

    // 仅返回必要字段，避免泄露多余信息
    result := make([]map[string]interface{}, 0, len(tenants))
    for _, t := range tenants {
        result = append(result, map[string]interface{}{
            "id":   t.ID,
            "code": t.Code,
            "name": t.Name,
        })
    }

    response.ReturnData(c, result)
}

