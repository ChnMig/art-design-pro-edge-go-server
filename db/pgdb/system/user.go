package system

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
	"api-server/util/encryption"
)

func encryptionPWD(password string) string {
	return encryption.MD5WithSalt(config.PWDSalt + password)
}

// 查询用户
func VerifyUser(userName, password string) (User, error) {
	user := User{}
	err := pgdb.GetClient().Where(&User{Username: userName, Password: encryptionPWD(password)}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		zap.L().Error("failed to get user", zap.Error(err))
		return user, err
	}
	return user, nil
}

func GetUser(user *User) error {
	if err := pgdb.GetClient().Where(user).First(user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return err
	}
	return nil
}
