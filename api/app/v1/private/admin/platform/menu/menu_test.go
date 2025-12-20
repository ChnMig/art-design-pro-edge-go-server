package menu

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-server/api/response"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUpdateTenantMenu_NotSuperAdmin(t *testing.T) {
	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(2))
		c.Next()
	})
	router.PUT("/tenant", UpdateTenantMenu)

	req, _ := http.NewRequest(http.MethodPut, "/tenant", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Code != response.PERMISSION_DENIED.Code {
		t.Errorf("Code = %d, want %d", resp.Code, response.PERMISSION_DENIED.Code)
	}
	if resp.Status != response.PERMISSION_DENIED.Status {
		t.Errorf("Status = %s, want %s", resp.Status, response.PERMISSION_DENIED.Status)
	}
	if resp.Message != "仅平台管理员可以调整租户菜单范围" {
		t.Errorf("Message = %q, want %q", resp.Message, "仅平台管理员可以调整租户菜单范围")
	}
}

