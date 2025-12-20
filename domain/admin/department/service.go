package department

import (
	"errors"

	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

type FindListQuery struct {
	Name   string
	Status uint
}

func FindDepartmentList(query FindListQuery, page, pageSize int) ([]system.SystemDepartment, int64, error) {
	filter := system.SystemDepartment{
		Name:   query.Name,
		Status: query.Status,
	}
	return system.FindDepartmentList(&filter, page, pageSize)
}

type AddInput struct {
	Name   string
	Status uint
	Sort   uint
}

func AddDepartment(input AddInput) (system.SystemDepartment, error) {
	department := system.SystemDepartment{
		Name:   input.Name,
		Status: input.Status,
		Sort:   input.Sort,
	}
	if err := system.AddDepartment(&department); err != nil {
		return system.SystemDepartment{}, err
	}
	return department, nil
}

type UpdateInput struct {
	ID     uint
	Name   string
	Status uint
	Sort   uint
}

func UpdateDepartment(input UpdateInput) (system.SystemDepartment, error) {
	department := system.SystemDepartment{
		Model:  gorm.Model{ID: input.ID},
		Name:   input.Name,
		Status: input.Status,
		Sort:   input.Sort,
	}
	if err := system.UpdateDepartment(&department); err != nil {
		return system.SystemDepartment{}, err
	}
	return department, nil
}

func DeleteDepartment(id uint) (system.SystemDepartment, error) {
	department := system.SystemDepartment{Model: gorm.Model{ID: id}}
	if err := system.GetDepartment(&department); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemDepartment{}, ErrDepartmentNotFound
		}
		return system.SystemDepartment{}, err
	}

	var userCount int64
	if err := system.CountUsersByDepartmentID(id, &userCount); err != nil {
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
