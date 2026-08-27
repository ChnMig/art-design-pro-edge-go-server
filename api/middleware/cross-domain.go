package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CorssDomainHandler consent cross-domain middleware
func CorssDomainHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               // method
		origin := c.Request.Header.Get("Origin") // header
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-ID")
			c.Header("Access-Control-Expose-Headers", "X-Request-ID")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}

		// Release all OPTIONS methods
		if method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		// Processing request
		c.Next()
	}
}
