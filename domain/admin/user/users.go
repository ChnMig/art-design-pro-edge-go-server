package user

import (
	"errors"

	"api-server/db/pgdb/system"
	systemuser "api-server/db/rdb/systemUser"

	"gorm.io/gorm"
)

type FindUserQuery struct {
	Username     string
	Name         string
	Phone        string
	DepartmentID uint
	RoleID       uint
}

func FindUserList(tenantID uint, query FindUserQuery, page, pageSize int) ([]system.UserWithRelations, int64, error) {
	filter := system.SystemUser{
		TenantID:     tenantID,
		Username:     query.Username,
		Name:         query.Name,
		Phone:        query.Phone,
		DepartmentID: query.DepartmentID,
		RoleID:       query.RoleID,
	}
	return system.FindUserList(&filter, page, pageSize)
}

type AddUserInput struct {
	Name         string
	Username     string
	Account      string
	Password     string
	Phone        string
	Gender       uint
	Status       uint
	RoleID       uint
	DepartmentID uint
}

func AddUser(tenantID uint, input AddUserInput) error {
	if tenantID == 0 {
		return ErrInvalidTenantContext
	}
	if len(input.Password) < 8 || len(input.Password) > 72 ||
		(input.Status != system.StatusEnabled && input.Status != system.StatusDisabled) {
		return ErrInvalidUserInput
	}
	if err := validateTenantAssignments(tenantID, input.RoleID, input.DepartmentID); err != nil {
		return err
	}

	u := system.SystemUser{
		TenantID:     tenantID,
		Name:         input.Name,
		Username:     input.Username,
		Account:      input.Account,
		Password:     input.Password,
		Phone:        input.Phone,
		Gender:       input.Gender,
		Status:       input.Status,
		RoleID:       input.RoleID,
		DepartmentID: input.DepartmentID,
	}
	if err := system.AddUser(&u); err != nil {
		return err
	}
	_ = systemuser.InvalidateUserCache(0)
	return nil
}

type UpdateUserInput struct {
	ID           uint
	Name         string
	Username     string
	Account      string
	Password     string
	Phone        string
	Gender       uint
	Status       uint
	RoleID       uint
	DepartmentID uint
}

func UpdateUser(tenantID uint, input UpdateUserInput) error {
	if tenantID == 0 {
		return ErrInvalidTenantContext
	}
	if input.ID == 1 && tenantID == system.DefaultTenantID {
		return ErrCannotModifySuperAdmin
	}
	if (input.Password != "" && (len(input.Password) < 8 || len(input.Password) > 72)) ||
		(input.Status != system.StatusEnabled && input.Status != system.StatusDisabled) {
		return ErrInvalidUserInput
	}
	if err := validateTenantAssignments(tenantID, input.RoleID, input.DepartmentID); err != nil {
		return err
	}

	u := system.SystemUser{
		Model:        gorm.Model{ID: input.ID},
		TenantID:     tenantID,
		Name:         input.Name,
		Username:     input.Username,
		Account:      input.Account,
		Phone:        input.Phone,
		Gender:       input.Gender,
		Status:       input.Status,
		RoleID:       input.RoleID,
		DepartmentID: input.DepartmentID,
	}
	if input.Password != "" {
		u.Password = input.Password
	}
	fields := []string{"name", "username", "account", "phone", "gender", "status", "role_id", "department_id"}
	if input.Password != "" {
		fields = append(fields, "password")
	}
	if err := system.UpdateUser(&u, fields...); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotInTenant
		}
		return err
	}
	_ = systemuser.InvalidateUserCache(input.ID)
	return nil
}

func validateTenantAssignments(tenantID, roleID, departmentID uint) error {
	roleEntity := system.SystemRole{Model: gorm.Model{ID: roleID}, TenantID: tenantID}
	if err := system.GetRole(&roleEntity); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotInTenant
		}
		return err
	}
	if roleEntity.Status != system.StatusEnabled {
		return ErrRoleNotInTenant
	}

	departmentEntity := system.SystemDepartment{Model: gorm.Model{ID: departmentID}, TenantID: tenantID}
	if err := system.GetDepartment(&departmentEntity); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDepartmentNotInTenant
		}
		return err
	}
	if departmentEntity.Status != system.StatusEnabled {
		return ErrDepartmentNotInTenant
	}
	return nil
}

func DeleteUser(tenantID, id uint) error {
	if tenantID == 0 {
		return ErrInvalidTenantContext
	}
	if id == 1 {
		return ErrCannotDeleteSuperAdmin
	}
	u := system.SystemUser{
		Model:    gorm.Model{ID: id},
		TenantID: tenantID,
	}
	if err := system.DeleteUser(&u); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotInTenant
		}
		return err
	}
	_ = systemuser.InvalidateUserCache(id)
	return nil
}

func IsRoleNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
