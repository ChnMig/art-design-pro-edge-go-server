package encryption

import (
	"testing"
)

func TestHashPasswordWithBcrypt(t *testing.T) {
	password := "testpassword123"

	// 测试密码哈希
	hashedPassword, err := HashPasswordWithBcrypt(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 检查哈希结果不为空
	if hashedPassword == "" {
		t.Fatal("Hashed password should not be empty")
	}

	// 检查哈希结果不等于原密码
	if hashedPassword == password {
		t.Fatal("Hashed password should not equal original password")
	}

	// 检查是否为bcrypt格式
	if !IsBcryptHash(hashedPassword) {
		t.Fatal("Result should be in bcrypt format")
	}
}

func TestVerifyBcryptPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	// 生成哈希
	hashedPassword, err := HashPasswordWithBcrypt(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 测试正确密码验证
	if !VerifyBcryptPassword(password, hashedPassword) {
		t.Fatal("Should verify correct password successfully")
	}

	// 测试错误密码验证
	if VerifyBcryptPassword(wrongPassword, hashedPassword) {
		t.Fatal("Should not verify wrong password")
	}
}

func TestIsBcryptHash(t *testing.T) {
	tests := []struct {
		name   string
		hash   string
		expect bool
	}{
		{"valid bcrypt hash", "$2a$10$abcdefghijklmnopqrstuvwxyz", true},
		{"valid bcrypt hash 2b", "$2b$12$abcdefghijklmnopqrstuvwxyz", true},
		{"md5 hash", "5d41402abc4b2a76b9719d911017c592", false},
		{"short string", "abc", false},
		{"empty string", "", false},
		{"not bcrypt format", "abc$def$ghi", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsBcryptHash(test.hash)
			if result != test.expect {
				t.Fatalf("Expected %v for hash %s, got %v", test.expect, test.hash, result)
			}
		})
	}
}
