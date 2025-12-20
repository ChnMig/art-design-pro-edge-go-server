package menu

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestGetMenuListByRoleID_MissingRoleID(t *testing.T) {
	router := setupTestRouter()
	router.GET("/menu", GetMenuListByRoleID)

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
	if resp.Code != response.INVALID_ARGUMENT.Code {
		t.Errorf("Code = %d, want %d", resp.Code, response.INVALID_ARGUMENT.Code)
	}
	if resp.Status != response.INVALID_ARGUMENT.Status {
		t.Errorf("Status = %s, want %s", resp.Status, response.INVALID_ARGUMENT.Status)
	}
	if resp.Message == "" {
		t.Errorf("Message should not be empty")
	}
}

func TestUpdateMenuListByRoleID_InvalidMenuData(t *testing.T) {
	router := setupTestRouter()
	router.POST("/menu", UpdateMenuListByRoleID)

	body := strings.NewReader(`{"role_id":1,"menu_data":"not-json"}`)
	req, _ := http.NewRequest(http.MethodPost, "/menu", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Code != response.DATA_LOSS.Code {
		t.Errorf("Code = %d, want %d", resp.Code, response.DATA_LOSS.Code)
	}
	if resp.Status != response.DATA_LOSS.Status {
		t.Errorf("Status = %s, want %s", resp.Status, response.DATA_LOSS.Status)
	}
	if resp.Message != "参数错误" {
		t.Errorf("Message = %q, want %q", resp.Message, "参数错误")
	}
}

