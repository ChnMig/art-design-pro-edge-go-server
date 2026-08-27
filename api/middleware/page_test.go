package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"api-server/config"
)

func newPageTestContext(rawQuery string) *gin.Context {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/?"+rawQuery, nil)
	return context
}

func TestGetPageRejectsNonPositiveValues(t *testing.T) {
	tests := []struct {
		name     string
		rawQuery string
		want     int
	}{
		{name: "正常页码", rawQuery: "page=3", want: 3},
		{name: "零", rawQuery: "page=0", want: config.DefaultPage},
		{name: "负数", rawQuery: "page=-1", want: config.DefaultPage},
		{name: "非数字", rawQuery: "page=abc", want: config.DefaultPage},
		{name: "缺省", rawQuery: "", want: config.DefaultPage},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GetPage(newPageTestContext(test.rawQuery)); got != test.want {
				t.Fatalf("GetPage() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestGetPageSizeEnforcesBounds(t *testing.T) {
	tests := []struct {
		name     string
		rawQuery string
		want     int
	}{
		{name: "下划线参数", rawQuery: "page_size=30", want: 30},
		{name: "驼峰参数", rawQuery: "pageSize=40", want: 40},
		{name: "零", rawQuery: "page_size=0", want: config.DefaultPageSize},
		{name: "负数", rawQuery: "page_size=-1", want: config.DefaultPageSize},
		{name: "超过上限", rawQuery: "page_size=101", want: config.MaxPageSize},
		{name: "非数字", rawQuery: "page_size=abc", want: config.DefaultPageSize},
		{name: "别名回退", rawQuery: "page_size=abc&pageSize=40", want: 40},
		{name: "缺省", rawQuery: "", want: config.DefaultPageSize},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GetPageSize(newPageTestContext(test.rawQuery)); got != test.want {
				t.Fatalf("GetPageSize() = %d, want %d", got, test.want)
			}
		})
	}
}
