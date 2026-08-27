package user

import (
	"errors"
	"testing"
)

func TestUpdateUserManagementCannotModifyPlatformRoot(t *testing.T) {
	err := UpdateUser(1, UpdateUserInput{ID: 1})
	if !errors.Is(err, ErrCannotModifySuperAdmin) {
		t.Fatalf("UpdateUser() error = %v, want ErrCannotModifySuperAdmin", err)
	}
}

func TestAddUserRejectsWeakPasswordBeforeDatabaseAccess(t *testing.T) {
	err := AddUser(2, AddUserInput{Password: "short", Status: 1})
	if !errors.Is(err, ErrInvalidUserInput) {
		t.Fatalf("AddUser() error = %v, want ErrInvalidUserInput", err)
	}
}
