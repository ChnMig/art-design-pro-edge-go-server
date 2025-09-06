package system

import (
	"testing"

	"api-server/config"
	"api-server/util/encryption"
)

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	
	// 测试MD5密码验证（向后兼容）
	md5Hash := encryption.MD5WithSalt(config.PWDSalt + password)
	if !VerifyPassword(password, md5Hash, "md5") {
		t.Fatal("Should verify MD5 password successfully")
	}
	
	// 测试bcrypt密码验证
	bcryptHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password with bcrypt: %v", err)
	}
	
	if !VerifyPassword(password, bcryptHash, "bcrypt") {
		t.Fatal("Should verify bcrypt password successfully")
	}
	
	// 测试自动检测bcrypt格式
	if !VerifyPassword(password, bcryptHash, "") {
		t.Fatal("Should auto-detect and verify bcrypt password")
	}
	
	// 测试错误密码
	if VerifyPassword("wrongpassword", bcryptHash, "bcrypt") {
		t.Fatal("Should not verify wrong password")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// 检查哈希结果
	if hashedPassword == "" {
		t.Fatal("Hashed password should not be empty")
	}
	
	if hashedPassword == password {
		t.Fatal("Hashed password should not equal original password")
	}
	
	// 验证哈希结果
	if !VerifyPassword(password, hashedPassword, "bcrypt") {
		t.Fatal("Should be able to verify the hashed password")
	}
}