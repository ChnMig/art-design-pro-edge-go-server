package tenant

import (
	"errors"
	"strings"

	"api-server/db/pgdb/system"
	systemuser "api-server/db/rdb/systemUser"

	"gorm.io/gorm"
)

type FindListQuery struct {
	Code   string
	Name   string
	Status uint
}

func FindTenantList(query FindListQuery, page, pageSize int) ([]system.SystemTenant, int64, error) {
	filter := system.SystemTenant{
		Code:   query.Code,
		Name:   query.Name,
		Status: query.Status,
	}
	return system.FindTenantList(&filter, page, pageSize)
}

type AddTenantInput struct {
	Code          string
	Name          string
	Contact       string
	Phone         string
	Email         string
	Status        uint
	AdminName     string
	AdminAccount  string
	AdminPassword string
}

func AddTenant(input AddTenantInput) (system.SystemTenant, error) {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.AdminName = strings.TrimSpace(input.AdminName)
	input.AdminAccount = strings.TrimSpace(input.AdminAccount)
	if input.Code == "" || input.Name == "" || input.AdminName == "" || len(input.AdminAccount) < 3 || len(input.AdminAccount) > 64 ||
		len(input.AdminPassword) < 8 || len(input.AdminPassword) > 72 ||
		(input.Status != system.StatusEnabled && input.Status != system.StatusDisabled) {
		return system.SystemTenant{}, ErrInvalidTenantInput
	}
	tenant := system.SystemTenant{
		Code:    input.Code,
		Name:    input.Name,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Status:  input.Status,
	}

	if err := system.CreateTenantWithAdmin(system.TenantBootstrap{
		Tenant:        &tenant,
		AdminName:     input.AdminName,
		AdminAccount:  input.AdminAccount,
		AdminPassword: input.AdminPassword,
	}); err != nil {
		return system.SystemTenant{}, err
	}
	_ = systemuser.InvalidateUserCache(0)
	return tenant, nil
}

type UpdateTenantInput struct {
	ID      uint
	Code    string
	Name    string
	Contact string
	Phone   string
	Email   string
	Status  uint
}

func UpdateTenant(input UpdateTenantInput) error {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	if input.Code == "" || input.Name == "" ||
		(input.Status != system.StatusEnabled && input.Status != system.StatusDisabled) {
		return ErrInvalidTenantInput
	}
	existing := system.SystemTenant{Model: gorm.Model{ID: input.ID}}
	if err := system.GetTenant(&existing); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return err
	}
	if input.ID == system.DefaultTenantID && (input.Status != system.StatusEnabled || input.Code != existing.Code) {
		return ErrDefaultTenantProtected
	}
	tenant := system.SystemTenant{
		Model:   gorm.Model{ID: input.ID},
		Code:    input.Code,
		Name:    input.Name,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Status:  input.Status,
	}

	if err := system.UpdateTenant(&tenant); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return err
	}
	return nil
}

func DeleteTenant(id uint) error {
	if id == system.DefaultTenantID {
		return ErrDefaultTenantProtected
	}
	tenant := system.SystemTenant{Model: gorm.Model{ID: id}}
	if err := system.GetTenant(&tenant); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return err
	}
	if err := system.DeleteTenant(&tenant); err != nil {
		return err
	}
	_ = systemuser.InvalidateAllUserCaches()
	return nil
}
