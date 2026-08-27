package system

import (
	"regexp"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newTenantScopeDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:              "host=localhost user=test dbname=test sslmode=disable",
		WithoutReturning: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("创建模拟数据库失败: %v", err)
	}
	return database
}

func assertScopedByIDAndTenant(t *testing.T, sql string) {
	t.Helper()
	wherePattern := regexp.MustCompile(`(?i)where[\s\S]*id[\s\S]*tenant_id`)
	if !wherePattern.MatchString(sql) {
		t.Fatalf("SQL 未同时按 ID 和租户限定: %s", sql)
	}
}

func TestUserMutationsAreTenantScoped(t *testing.T) {
	database := newTenantScopeDryRunDB(t)
	user := SystemUser{TenantID: 22, Name: "新名称"}
	user.ID = 11

	update := scopedUserUpdate(database, &user, []string{"name"})
	if update.Error != nil {
		t.Fatalf("构建用户更新 SQL 失败: %v", update.Error)
	}
	updateSQL := update.Statement.SQL.String()
	assertScopedByIDAndTenant(t, updateSQL)
	if strings.Contains(strings.ToLower(updateSQL), `set "tenant_id"`) {
		t.Fatalf("更新用户时不应修改租户归属: %s", updateSQL)
	}

	deleteResult := scopedUserDelete(database, &user)
	if deleteResult.Error != nil {
		t.Fatalf("构建用户删除 SQL 失败: %v", deleteResult.Error)
	}
	assertScopedByIDAndTenant(t, deleteResult.Statement.SQL.String())
}

func TestUserListRequiresTenant(t *testing.T) {
	database := newTenantScopeDryRunDB(t)
	query, err := buildUserListQuery(database, &SystemUser{TenantID: 22, Username: "alice"})
	if err != nil {
		t.Fatalf("buildUserListQuery() error = %v", err)
	}
	result := query.Find(&[]UserWithRelations{})
	if !strings.Contains(strings.ToLower(result.Statement.SQL.String()), "system_users.tenant_id =") {
		t.Fatalf("用户列表 SQL 缺少租户条件: %s", result.Statement.SQL.String())
	}

	if _, err := buildUserListQuery(database, &SystemUser{}); err == nil {
		t.Fatal("缺少租户标识时应拒绝构建用户列表查询")
	}
}

func TestDepartmentMutationsAreTenantScoped(t *testing.T) {
	database := newTenantScopeDryRunDB(t)
	department := SystemDepartment{TenantID: 22, Name: "新部门", Status: StatusEnabled}
	department.ID = 33

	update := scopedDepartmentUpdate(database, &department)
	if update.Error != nil {
		t.Fatalf("构建部门更新 SQL 失败: %v", update.Error)
	}
	updateSQL := update.Statement.SQL.String()
	assertScopedByIDAndTenant(t, updateSQL)
	if strings.Contains(strings.ToLower(updateSQL), `set "tenant_id"`) {
		t.Fatalf("更新部门时不应修改租户归属: %s", updateSQL)
	}

	deleteResult := scopedDepartmentDelete(database, &department)
	if deleteResult.Error != nil {
		t.Fatalf("构建部门删除 SQL 失败: %v", deleteResult.Error)
	}
	assertScopedByIDAndTenant(t, deleteResult.Statement.SQL.String())
}

func TestDepartmentListRequiresTenant(t *testing.T) {
	database := newTenantScopeDryRunDB(t)
	query, err := buildDepartmentListQuery(database, &SystemDepartment{TenantID: 22, Name: "研发"})
	if err != nil {
		t.Fatalf("buildDepartmentListQuery() error = %v", err)
	}
	result := query.Find(&[]SystemDepartment{})
	if !strings.Contains(strings.ToLower(result.Statement.SQL.String()), "tenant_id =") {
		t.Fatalf("部门列表 SQL 缺少租户条件: %s", result.Statement.SQL.String())
	}

	if _, err := buildDepartmentListQuery(database, &SystemDepartment{}); err == nil {
		t.Fatal("缺少租户标识时应拒绝构建部门列表查询")
	}
}

func TestLoginLogQueryRequiresExactTenantCode(t *testing.T) {
	database := newTenantScopeDryRunDB(t)
	query, err := buildLoginLogQuery(database, &SystemUserLoginLog{TenantCode: "tenant-a", UserName: "alice"})
	if err != nil {
		t.Fatalf("buildLoginLogQuery() error = %v", err)
	}
	result := query.Find(&[]SystemUserLoginLog{})
	sql := strings.ToLower(result.Statement.SQL.String())
	if !strings.Contains(sql, "tenant_code =") {
		t.Fatalf("登录日志 SQL 缺少精确租户条件: %s", sql)
	}

	if _, err := buildLoginLogQuery(database, &SystemUserLoginLog{}); err == nil {
		t.Fatal("缺少租户编码时应拒绝构建登录日志查询")
	}
}
