package middleware

import (
	"github.com/gin-gonic/gin"

	"api-server/api/response"
	"api-server/util/authentication"
)

const (
	JWTDataKey = "jwtData"
)

// TokenVerify Get the token and verify its validity
func TokenVerify(c *gin.Context) {
	c.FormFile("file") // 防止文件未发送完成就返回错误, 导致前端504而不是正确响应
	token := c.Request.Header.Get("Access-Token")
	if token == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	data, err := authentication.JWTDecrypt(token)
	if err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "token 解析失败")
		return
	}
	// set data to gin.Context
	c.Set(JWTDataKey, data)
	// Next
	c.Next()
}
