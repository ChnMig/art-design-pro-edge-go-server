package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/auth"
	"api-server/config"
	"api-server/db/pgdb/system"
)

func issueTestToken(t *testing.T, userID, tenantID uint, account string) string {
	t.Helper()
	oldKey := config.JWTKey
	oldExpiration := config.JWTExpiration
	config.JWTKey = "middleware-test-key-with-sufficient-length"
	config.JWTExpiration = time.Hour
	t.Cleanup(func() {
		config.JWTKey = oldKey
		config.JWTExpiration = oldExpiration
	})

	token, err := auth.JWTIssue(userID, tenantID, account)
	if err != nil {
		t.Fatalf("签发测试令牌失败: %v", err)
	}
	return token
}

func runTokenVerifyForTest(token string, loader tokenPrincipalLoader) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)
	tokenVerify(context, loader)
	return context, recorder
}

func TestTokenVerifyRechecksPrincipalState(t *testing.T) {
	token := issueTestToken(t, 7, 9, "alice")
	enabledUser := system.SystemUser{TenantID: 9, Account: "alice", Status: system.StatusEnabled, Password: "secret"}
	enabledUser.ID = 7
	enabledTenant := system.SystemTenant{Status: system.StatusEnabled}
	enabledTenant.ID = 9

	tests := []struct {
		name       string
		user       system.SystemUser
		tenant     system.SystemTenant
		loadErr    error
		wantDenied bool
	}{
		{name: "有效身份", user: enabledUser, tenant: enabledTenant},
		{name: "用户已删除", user: enabledUser, tenant: enabledTenant, loadErr: gorm.ErrRecordNotFound, wantDenied: true},
		{name: "用户已禁用", user: func() system.SystemUser { value := enabledUser; value.Status = system.StatusDisabled; return value }(), tenant: enabledTenant, wantDenied: true},
		{name: "租户已禁用", user: enabledUser, tenant: func() system.SystemTenant { value := enabledTenant; value.Status = system.StatusDisabled; return value }(), wantDenied: true},
		{name: "用户已迁出租户", user: func() system.SystemUser { value := enabledUser; value.TenantID = 10; return value }(), tenant: enabledTenant, wantDenied: true},
		{name: "账号已变更", user: func() system.SystemUser { value := enabledUser; value.Account = "renamed"; return value }(), tenant: enabledTenant, wantDenied: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, recorder := runTokenVerifyForTest(token, func(_, _ uint) (system.SystemUser, system.SystemTenant, error) {
				return test.user, test.tenant, test.loadErr
			})

			if context.IsAborted() != test.wantDenied {
				t.Fatalf("IsAborted() = %v, want %v", context.IsAborted(), test.wantDenied)
			}
			if test.wantDenied {
				if !strings.Contains(recorder.Body.String(), `"code":401`) {
					t.Fatalf("拒绝响应未返回 401 业务码: %s", recorder.Body.String())
				}
				return
			}

			currentUser, ok := GetCurrentUser(context)
			if !ok || currentUser.ID != enabledUser.ID || currentUser.TenantID != enabledUser.TenantID {
				t.Fatalf("当前用户上下文错误: %#v, ok=%v", currentUser, ok)
			}
			if currentUser.Password != "" {
				t.Fatal("当前用户上下文不应保留密码")
			}
		})
	}
}

func TestTokenVerifyRejectsLoaderFailure(t *testing.T) {
	token := issueTestToken(t, 3, 4, "bob")
	context, _ := runTokenVerifyForTest(token, func(_, _ uint) (system.SystemUser, system.SystemTenant, error) {
		return system.SystemUser{}, system.SystemTenant{}, errors.New("数据库不可用")
	})
	if !context.IsAborted() {
		t.Fatal("身份查询失败时应拒绝请求")
	}
}
