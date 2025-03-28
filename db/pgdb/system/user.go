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

func FindUserList(user *User, page, pageSize int) ([]User, error) {
	var users []User
	db := pgdb.GetClient()
	// 构建查询条件
	query := db
	// 使用模糊查询
	if user.Username != "" {
		query = query.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user.Name != "" {
		query = query.Where("name LIKE ?", "%"+user.Name+"%")
	}
	if user.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+user.Phone+"%")
	}
	// 应用分页
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		zap.L().Error("failed to find user list", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func AddUser(user *User) error {
	if err := pgdb.GetClient().Create(user).Error; err != nil {
		zap.L().Error("failed to add user", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUser(user *User) error {
	if err := pgdb.GetClient().Save(user).Error; err != nil {
		zap.L().Error("failed to update user", zap.Error(err))
		return err
	}
	return nil
}

func DeleteUser(user *User) error {
	if err := pgdb.GetClient().Delete(user).Error; err != nil {
		zap.L().Error("failed to delete user", zap.Error(err))
		return err
	}
	return nil
}
