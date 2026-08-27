package role

import (
	"errors"
	"strings"

	"api-server/db/pgdb/system"
	systemuser "api-server/db/rdb/systemUser"

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
	TenantID      uint
	Name          string
	Status        uint
	Desc          string
	IsTenantAdmin bool
}

func AddRole(input AddInput) (system.SystemRole, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.TenantID == 0 || input.Name == "" ||
		(input.Status != system.StatusEnabled && input.Status != system.StatusDisabled) {
		return system.SystemRole{}, ErrInvalidRoleInput
	}
	tenant := system.SystemTenant{Model: gorm.Model{ID: input.TenantID}}
	if err := system.GetTenant(&tenant); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemRole{}, ErrRoleTenantNotFound
		}
		return system.SystemRole{}, err
	}
	role := system.SystemRole{
		TenantID:      input.TenantID,
		Name:          input.Name,
		Status:        input.Status,
		Desc:          input.Desc,
		IsTenantAdmin: input.IsTenantAdmin,
	}
	if err := system.AddRole(&role); err != nil {
		return system.SystemRole{}, err
	}
	return role, nil
}

type UpdateInput struct {
	ID            uint
	TenantID      uint
	Name          string
	Status        uint
	Desc          string
	IsTenantAdmin *bool
}

func UpdateRole(input UpdateInput) (system.SystemRole, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.ID == 0 || input.Name == "" ||
		(input.Status != system.StatusEnabled && input.Status != system.StatusDisabled) {
		return system.SystemRole{}, ErrInvalidRoleInput
	}
	existing, err := GetRole(input.ID)
	if err != nil {
		return system.SystemRole{}, err
	}

	targetTenantID := existing.TenantID
	if input.TenantID != 0 && input.TenantID != existing.TenantID {
		return system.SystemRole{}, ErrRoleTenantImmutable
	}
	isTenantAdmin := existing.IsTenantAdmin
	if input.IsTenantAdmin != nil {
		isTenantAdmin = *input.IsTenantAdmin
	}
	if input.Status != system.StatusEnabled || (existing.IsTenantAdmin && !isTenantAdmin) {
		count, countErr := system.CountUsersByRoleID(existing.TenantID, existing.ID)
		if countErr != nil {
			return system.SystemRole{}, countErr
		}
		if count > 0 {
			return system.SystemRole{}, ErrRoleInUse
		}
	}

	role := system.SystemRole{
		Model:         gorm.Model{ID: input.ID},
		TenantID:      targetTenantID,
		Name:          input.Name,
		Status:        input.Status,
		Desc:          input.Desc,
		IsTenantAdmin: isTenantAdmin,
	}
	if err := system.UpdateRole(&role); err != nil {
		return system.SystemRole{}, err
	}
	_ = systemuser.InvalidateAllUserCaches()
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
	if role.IsTenantAdmin {
		return ErrCannotDeleteTenantAdmin
	}
	count, err := system.CountUsersByRoleID(role.TenantID, role.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleInUse
	}
	if err := system.DeleteRole(&role); err != nil {
		return err
	}
	_ = systemuser.InvalidateAllUserCaches()
	return nil
}
