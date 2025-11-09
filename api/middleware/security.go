package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-server/api/response"
)

// BodySizeLimit 限制请求体大小
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()

		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			response.ReturnError(c, response.INVALID_ARGUMENT, "请求体过大")
			return
		}
	}
}

// SecurityHeaders 添加常见安全响应头
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
