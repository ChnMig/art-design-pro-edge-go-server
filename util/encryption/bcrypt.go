package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// BCryptCost bcrypt算法的成本因子，12提供良好的安全性/性能平衡
	BCryptCost = 12
)

// HashPasswordWithBcrypt 使用bcrypt算法对密码进行哈希
func HashPasswordWithBcrypt(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BCryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyBcryptPassword 验证bcrypt哈希的密码
func VerifyBcryptPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// IsBcryptHash 检查是否为bcrypt哈希格式
// bcrypt哈希格式: $2a$10$... 或 $2b$10$... 等
func IsBcryptHash(hash string) bool {
	if len(hash) < 7 {
		return false
	}
	return hash[0] == '$' && (hash[1] == '2') && (hash[2] == 'a' || hash[2] == 'b' || hash[2] == 'x' || hash[2] == 'y')
}
