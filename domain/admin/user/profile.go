package user

import (
	"errors"

	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

type UpdateProfileInput struct {
	UserID   uint
	Password string
	Username string
	Phone    string
	Gender   uint
}

func UpdateUserProfile(input UpdateProfileInput) error {
	u := system.SystemUser{
		Model:    gorm.Model{ID: input.UserID},
		Username: input.Username,
		Phone:    input.Phone,
		Gender:   input.Gender,
	}
	if input.Password != "" {
		u.Password = input.Password
	}
	return system.UpdateUser(&u)
}

func GetUserProfile(userID uint) (system.SystemUser, error) {
	user := system.SystemUser{Model: gorm.Model{ID: userID}}
	if err := system.GetUser(&user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return system.SystemUser{}, ErrUserNotFound
		}
		return system.SystemUser{}, err
	}
	user.Password = ""
	return user, nil
}

