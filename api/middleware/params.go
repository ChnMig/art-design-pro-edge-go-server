package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"api-server/api/response"
)

// 校验参数
func CheckParam(params interface{}, c *gin.Context) bool {
	if err := c.ShouldBindWith(params, binding.Default(c.Request.Method, c.ContentType())); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return false
	}
	return true
}
