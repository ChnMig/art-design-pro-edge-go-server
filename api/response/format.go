package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ReturnErrorWithData(c *gin.Context, data responseData, result interface{}) {
	data.Timestamp = time.Now().Unix()
	data.Data = result
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseOk 正常响应
func ReturnOk(c *gin.Context, result interface{}) {
	data := OK
	data.Timestamp = time.Now().Unix()
	data.Data = result
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseOkWithCount 带数量的正常响应
func ReturnOkWithCount(c *gin.Context, count int, result interface{}) {
	data := OK
	data.Timestamp = time.Now().Unix()
	data.Data = result
	data.Count = &count
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseError 错误响应
func ReturnError(c *gin.Context, data responseData, description string) {
	data.Timestamp = time.Now().Unix()
	data.Message = func() string {
		if description == "" {
			return data.Message
		}
		return description
	}()
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseSuccess 执行成功
func ReturnSuccess(c *gin.Context) {
	data := OK
	data.Timestamp = time.Now().Unix()
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}
