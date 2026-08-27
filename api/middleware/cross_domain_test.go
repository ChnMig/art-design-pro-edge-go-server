package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCrossDomainOptionsStopsHandlerChain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	called := false
	router := gin.New()
	router.Use(CorssDomainHandler())
	router.OPTIONS("/resource", func(c *gin.Context) {
		called = true
		c.Status(http.StatusInternalServerError)
	})

	request := httptest.NewRequest(http.MethodOptions, "/resource", nil)
	request.Header.Set("Origin", "https://example.test")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if called {
		t.Fatal("OPTIONS 预检请求不应继续执行后续处理器")
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS 状态码 = %d，期望 %d", recorder.Code, http.StatusNoContent)
	}
	if recorder.Header().Get("Access-Control-Allow-Headers") != "Authorization,Content-Type,X-Request-ID" {
		t.Fatalf("允许请求头错误: %q", recorder.Header().Get("Access-Control-Allow-Headers"))
	}
}
