package tenant

import (
	"errors"

	"api-server/db/pgdb/system"

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
	Code    string
	Name    string
	Contact string
	Phone   string
	Email   string
	Status  uint
}

func AddTenant(input AddTenantInput) (system.SystemTenant, error) {
	tenant := system.SystemTenant{
		Code:    input.Code,
		Name:    input.Name,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Status:  input.Status,
	}

	if err := system.AddTenant(&tenant); err != nil {
		return system.SystemTenant{}, err
	}
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
		return err
	}
	return nil
}

func DeleteTenant(id uint) error {
	tenant := system.SystemTenant{Model: gorm.Model{ID: id}}
	if err := system.GetTenant(&tenant); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return err
	}
	return system.DeleteTenant(&tenant)
}

