package authentication

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"api-server/config"
)

func withJWTTestConfig(t *testing.T) {
	t.Helper()
	oldKey := config.JWTKey
	oldExpiration := config.JWTExpiration
	config.JWTKey = "jwt-validation-test-key-at-least-32-characters"
	config.JWTExpiration = time.Hour
	t.Cleanup(func() {
		config.JWTKey = oldKey
		config.JWTExpiration = oldExpiration
	})
}

func TestParseHS256RejectsOtherHMACAlgorithms(t *testing.T) {
	withJWTTestConfig(t)
	claims := &MapClaims{Data: map[string]interface{}{"user_id": 1}}
	PrepareRegisteredClaims(&claims.RegisteredClaims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte(config.JWTKey))
	if err != nil {
		t.Fatalf("签发 HS512 测试令牌失败: %v", err)
	}
	if _, err := ParseHS256(signed, &MapClaims{}); err == nil {
		t.Fatal("ParseHS256() 不应接受 HS512 令牌")
	}
}

func TestParseHS256RequiresRegisteredClaims(t *testing.T) {
	withJWTTestConfig(t)
	claims := &MapClaims{
		Data: map[string]interface{}{"user_id": 1},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   defaultIssuer,
			Subject:  defaultSubject,
			Audience: jwt.ClaimStrings{defaultAudience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JWTKey))
	if err != nil {
		t.Fatalf("签发缺少过期时间的测试令牌失败: %v", err)
	}
	if _, err := ParseHS256(signed, &MapClaims{}); err == nil {
		t.Fatal("ParseHS256() 不应接受缺少 exp 的令牌")
	}
}
