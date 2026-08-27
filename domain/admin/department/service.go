package department

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

func FindDepartmentList(query FindListQuery, page, pageSize int) ([]system.SystemDepartment, int64, error) {
	filter := system.SystemDepartment{
		TenantID: query.TenantID,
		Name:     query.Name,
		Status:   query.Status,
	}
	return system.FindDepartmentList(&filter, page, pageSize)
}

type AddInput struct {
	TenantID uint
	Name     string
	Status   uint
	Sort     uint
}

func AddDepartment(input AddInput) (system.SystemDepartment, error) {
	department := system.SystemDepartment{
		TenantID: input.TenantID,
		Name:     input.Name,
		Status:   input.Status,
		Sort:     input.Sort,
	}
	if err := system.AddDepartment(&department); err != nil {
		return system.SystemDepartment{}, err
	}
	return department, nil
}

type UpdateInput struct {
	TenantID uint
	ID       uint
	Name     string
	Status   uint
	Sort     uint
}

func UpdateDepartment(input UpdateInput) (system.SystemDepartment, error) {
	department := system.SystemDepartment{
		Model:    gorm.Model{ID: input.ID},
		TenantID: input.TenantID,
		Name:     input.Name,
		Status:   input.Status,
		Sort:     input.Sort,
	}
	if err := system.UpdateDepartment(&department); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemDepartment{}, ErrDepartmentNotFound
		}
		return system.SystemDepartment{}, err
	}
	return department, nil
}

func DeleteDepartment(tenantID, id uint) (system.SystemDepartment, error) {
	department := system.SystemDepartment{Model: gorm.Model{ID: id}, TenantID: tenantID}
	if err := system.GetDepartment(&department); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemDepartment{}, ErrDepartmentNotFound
		}
		return system.SystemDepartment{}, err
	}

	var userCount int64
	if err := system.CountUsersByDepartmentID(tenantID, id, &userCount); err != nil {
		return system.SystemDepartment{}, err
	}
	if userCount > 0 {
		return system.SystemDepartment{}, ErrDepartmentHasUsers
	}

	if err := system.DeleteDepartment(&department); err != nil {
		return system.SystemDepartment{}, err
	}
	return department, nil
}
