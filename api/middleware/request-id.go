package middleware

import (
	"api-server/util/id"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const RequestIDKey = "X-Request-ID"

// RequestID 为请求生成唯一 ID，并注入上下文 logger
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = id.GenerateID()
		}

		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)

		contextLogger := zap.L().With(
			zap.String("trace_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)
		c.Set("logger", contextLogger)

		contextLogger.Info("请求开始")
		c.Next()
		contextLogger.Info("请求结束", zap.Int("status_code", c.Writer.Status()))
	}
}
