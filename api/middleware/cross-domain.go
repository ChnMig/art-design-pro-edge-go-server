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
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")  // This is to allow access to all domains
			c.Header("Access-Control-Allow-Methods", "*") // All cross-domain request methods supported by the server, in order to avoid multiple'pre-check' requests for browsing requests
			//  header
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Expose-Headers", "*")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}

		// Release all OPTIONS methods
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// Processing request
		c.Next()
	}
}
