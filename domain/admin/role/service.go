package role

import (
	"errors"

	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

type FindListQuery struct {
	TenantID uint
	Name     string
	Status   uint
}

func GetRole(id uint) (system.SystemRole, error) {
	role := system.SystemRole{Model: gorm.Model{ID: id}}
	if err := system.GetRole(&role); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemRole{}, ErrRoleNotFound
		}
		return system.SystemRole{}, err
	}
	return role, nil
}

func FindRoleList(query FindListQuery, page, pageSize int) ([]system.SystemRole, int64, error) {
	filter := system.SystemRole{
		TenantID: query.TenantID,
		Name:     query.Name,
		Status:   query.Status,
	}
	return system.FindRoleList(&filter, page, pageSize)
}

type AddInput struct {
	TenantID uint
	Name     string
	Status   uint
	Desc     string
}

func AddRole(input AddInput) (system.SystemRole, error) {
	role := system.SystemRole{
		TenantID: input.TenantID,
		Name:     input.Name,
		Status:   input.Status,
		Desc:     input.Desc,
	}
	if err := system.AddRole(&role); err != nil {
		return system.SystemRole{}, err
	}
	return role, nil
}

type UpdateInput struct {
	ID       uint
	TenantID uint
	Name     string
	Status   uint
	Desc     string
}

func UpdateRole(input UpdateInput) (system.SystemRole, error) {
	existing, err := GetRole(input.ID)
	if err != nil {
		return system.SystemRole{}, err
	}

	targetTenantID := existing.TenantID
	if input.TenantID != 0 {
		targetTenantID = input.TenantID
	}

	role := system.SystemRole{
		Model:    gorm.Model{ID: input.ID},
		TenantID: targetTenantID,
		Name:     input.Name,
		Status:   input.Status,
		Desc:     input.Desc,
	}
	if err := system.UpdateRole(&role); err != nil {
		return system.SystemRole{}, err
	}
	return role, nil
}

func DeleteRole(id uint) error {
	role := system.SystemRole{Model: gorm.Model{ID: id}}
	if err := system.GetRole(&role); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}
	return system.DeleteRole(&role)
}
