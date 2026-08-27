package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-server/api/response"
)

func TestGetCaptchaRejectsUnsafeDimensions(t *testing.T) {
	router := setupTestRouter()
	router.GET("/captcha", GetCaptcha)

	request := httptest.NewRequest(http.MethodGet, "/captcha?width=1&height=1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	var result errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if result.Code != response.INVALID_ARGUMENT.Code {
		t.Fatalf("Code = %d, want %d", result.Code, response.INVALID_ARGUMENT.Code)
	}
	if result.Message != "验证码宽度需为 80-600，高度需为 40-200" {
		t.Fatalf("Message = %q", result.Message)
	}
}
