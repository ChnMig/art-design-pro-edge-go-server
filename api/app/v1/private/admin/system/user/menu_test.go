package user

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

func TestGetUserMenuList_MissingToken(t *testing.T) {
	router := setupTestRouter()
	router.GET("/menu", GetUserMenuList)

	req, _ := http.NewRequest(http.MethodGet, "/menu", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Code != response.UNAUTHENTICATED.Code {
		t.Errorf("Code = %d, want %d", resp.Code, response.UNAUTHENTICATED.Code)
	}
	if resp.Status != response.UNAUTHENTICATED.Status {
		t.Errorf("Status = %s, want %s", resp.Status, response.UNAUTHENTICATED.Status)
	}
	if resp.Message != "未携带 token" {
		t.Errorf("Message = %q, want %q", resp.Message, "未携带 token")
	}
}

func TestGetUserMenuList_MissingTenant(t *testing.T) {
	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(2))
		c.Next()
	})
	router.GET("/menu", GetUserMenuList)

	req, _ := http.NewRequest(http.MethodGet, "/menu", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Code != response.UNAUTHENTICATED.Code {
		t.Errorf("Code = %d, want %d", resp.Code, response.UNAUTHENTICATED.Code)
	}
	if resp.Status != response.UNAUTHENTICATED.Status {
		t.Errorf("Status = %s, want %s", resp.Status, response.UNAUTHENTICATED.Status)
	}
	if resp.Message != "租户信息缺失" {
		t.Errorf("Message = %q, want %q", resp.Message, "租户信息缺失")
	}
}

