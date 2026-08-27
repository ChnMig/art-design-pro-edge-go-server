package user

import (
	"errors"

	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

type UpdateProfileInput struct {
	TenantID uint
	UserID   uint
	Password string
	Username string
	Phone    string
	Gender   uint
}

func UpdateUserProfile(input UpdateProfileInput) error {
	u := system.SystemUser{
		Model:    gorm.Model{ID: input.UserID},
		TenantID: input.TenantID,
		Username: input.Username,
		Phone:    input.Phone,
		Gender:   input.Gender,
	}
	fields := []string{"username", "phone", "gender"}
	if input.Password != "" {
		u.Password = input.Password
		fields = append(fields, "password")
	}
	if err := system.UpdateUser(&u, fields...); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotInTenant
		}
		return err
	}
	return nil
}

func GetUserProfile(tenantID, userID uint) (system.SystemUser, error) {
	if tenantID == 0 || userID == 0 {
		return system.SystemUser{}, ErrInvalidTenantContext
	}
	user := system.SystemUser{Model: gorm.Model{ID: userID}, TenantID: tenantID}
	if err := system.GetUser(&user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemUser{}, ErrUserNotFound
		}
		return system.SystemUser{}, err
	}
	user.Password = ""
	return user, nil
}
