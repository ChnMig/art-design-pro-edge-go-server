package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

type statusResponse struct {
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   StatusDTO `json:"data"`
}

func TestStatus(t *testing.T) {
	router := setupTestRouter()
	router.GET("/health", Status)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status() status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp statusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("Status() response code = %d, want 200", resp.Code)
	}

	if resp.Status != "OK" {
		t.Errorf("Status() wrapper status = %s, want 'OK'", resp.Status)
	}

	if resp.Data.Status != "ok" {
		t.Errorf("Status() data.status = %s, want 'ok'", resp.Data.Status)
	}
	if !resp.Data.Ready {
		t.Errorf("Status() data.ready = %v, want true", resp.Data.Ready)
	}
	if resp.Data.Uptime == "" {
		t.Errorf("Status() missing uptime field")
	}
	if resp.Data.Timestamp == 0 {
		t.Errorf("Status() missing timestamp field")
	}
}
