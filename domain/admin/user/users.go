package user

import (
	"errors"

	"api-server/db/pgdb/system"

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
	roleEntity := system.SystemRole{Model: gorm.Model{ID: input.RoleID}}
	if err := system.GetRole(&roleEntity); err != nil || roleEntity.TenantID != tenantID {
		return ErrRoleNotInTenant
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
	return system.AddUser(&u)
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
	roleEntity := system.SystemRole{Model: gorm.Model{ID: input.RoleID}}
	if err := system.GetRole(&roleEntity); err != nil || roleEntity.TenantID != tenantID {
		return ErrRoleNotInTenant
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
	return system.UpdateUser(&u)
}

func DeleteUser(id uint) error {
	if id == 1 {
		return ErrCannotDeleteSuperAdmin
	}
	u := system.SystemUser{
		Model: gorm.Model{ID: id},
	}
	if err := system.DeleteUser(&u); err != nil {
		return err
	}
	return nil
}

func IsRoleNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

