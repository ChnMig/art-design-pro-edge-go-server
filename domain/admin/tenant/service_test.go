package tenant

import (
	"errors"
	"testing"

	"api-server/db/pgdb/system"
)

func TestDeleteTenantProtectsDefaultTenantBeforeDatabaseAccess(t *testing.T) {
	if err := DeleteTenant(system.DefaultTenantID); !errors.Is(err, ErrDefaultTenantProtected) {
		t.Fatalf("DeleteTenant() error = %v, want ErrDefaultTenantProtected", err)
	}
}

func TestAddTenantRejectsWeakBootstrapPassword(t *testing.T) {
	_, err := AddTenant(AddTenantInput{
		Code: "tenant-a", Name: "租户甲", Status: system.StatusEnabled,
		AdminName: "管理员", AdminAccount: "admin", AdminPassword: "short",
	})
	if !errors.Is(err, ErrInvalidTenantInput) {
		t.Fatalf("AddTenant() error = %v, want ErrInvalidTenantInput", err)
	}
}
