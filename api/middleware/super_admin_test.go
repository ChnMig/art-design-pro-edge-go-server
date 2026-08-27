package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/db/pgdb/system"
)

func newAdminTestContext(user *system.SystemUser, tenantID uint) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	context.Set("tenant_id", tenantID)
	tenant := system.SystemTenant{Status: system.StatusEnabled}
	tenant.ID = tenantID
	context.Set(currentTenantContextKey, tenant)
	if user != nil {
		if user.Status == 0 {
			user.Status = system.StatusEnabled
		}
		context.Set(currentUserContextKey, *user)
	}
	return context, recorder
}

func TestTenantAdminVerify(t *testing.T) {
	tests := []struct {
		name       string
		user       *system.SystemUser
		tenantID   uint
		role       *system.SystemRole
		wantDenied bool
		wantCode   string
	}{
		{name: "超级管理员", user: &system.SystemUser{Model: modelWithID(1), TenantID: 1}, tenantID: 1},
		{name: "租户管理员", user: &system.SystemUser{Model: modelWithID(8), TenantID: 3, RoleID: 27}, tenantID: 3, role: &system.SystemRole{Model: modelWithID(27), TenantID: 3, Status: system.StatusEnabled, IsTenantAdmin: true}},
		{name: "固定角色ID不再授权", user: &system.SystemUser{Model: modelWithID(9), TenantID: 3, RoleID: 2}, tenantID: 3, role: &system.SystemRole{Model: modelWithID(2), TenantID: 3, Status: system.StatusEnabled}, wantDenied: true, wantCode: `"code":403`},
		{name: "禁用的管理员角色", user: &system.SystemUser{Model: modelWithID(9), TenantID: 3, RoleID: 27}, tenantID: 3, role: &system.SystemRole{Model: modelWithID(27), TenantID: 3, Status: system.StatusDisabled, IsTenantAdmin: true}, wantDenied: true, wantCode: `"code":403`},
		{name: "租户不匹配", user: &system.SystemUser{Model: modelWithID(8), TenantID: 3, RoleID: 27}, tenantID: 4, wantDenied: true, wantCode: `"code":401`},
		{name: "缺少已验证用户", tenantID: 3, wantDenied: true, wantCode: `"code":401`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, recorder := newAdminTestContext(test.user, test.tenantID)
			tenantAdminVerify(context, func(role *system.SystemRole) error {
				if test.role == nil {
					return gorm.ErrRecordNotFound
				}
				*role = *test.role
				return nil
			})
			if context.IsAborted() != test.wantDenied {
				t.Fatalf("IsAborted() = %v, want %v", context.IsAborted(), test.wantDenied)
			}
			if test.wantCode != "" && !strings.Contains(recorder.Body.String(), test.wantCode) {
				t.Fatalf("响应业务码错误: %s", recorder.Body.String())
			}
		})
	}
}

func TestSuperAdminVerify(t *testing.T) {
	tests := []struct {
		name       string
		user       *system.SystemUser
		wantDenied bool
		wantCode   string
	}{
		{name: "超级管理员", user: &system.SystemUser{Model: modelWithID(1), TenantID: 1}},
		{name: "其他租户的用户1", user: &system.SystemUser{Model: modelWithID(1), TenantID: 2}, wantDenied: true, wantCode: `"code":401`},
		{name: "非超级管理员", user: &system.SystemUser{Model: modelWithID(2), TenantID: 1}, wantDenied: true, wantCode: `"code":403`},
		{name: "缺少已验证用户", wantDenied: true, wantCode: `"code":401`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, recorder := newAdminTestContext(test.user, 1)
			SuperAdminVerify(context)
			if context.IsAborted() != test.wantDenied {
				t.Fatalf("IsAborted() = %v, want %v", context.IsAborted(), test.wantDenied)
			}
			if test.wantCode != "" && !strings.Contains(recorder.Body.String(), test.wantCode) {
				t.Fatalf("响应业务码错误: %s", recorder.Body.String())
			}
		})
	}
}

func modelWithID(id uint) gorm.Model {
	return gorm.Model{ID: id}
}
