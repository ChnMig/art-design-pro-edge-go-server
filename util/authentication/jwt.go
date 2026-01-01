package authentication

import (
	"fmt"
	"time"

	"api-server/config"
	"api-server/util/id"

	"github.com/golang-jwt/jwt/v5"
)

var (
	defaultIssuer   = "http-services"
	defaultSubject  = "token"
	defaultAudience = "client"
)

// MapClaims 灵活的 Claims，使用 map 存储自定义数据
// 适用于不同项目有不同数据结构的场景
type MapClaims struct {
	Data map[string]interface{} `json:"data"` // 自定义数据，完全灵活
	jwt.RegisteredClaims
}

// PrepareRegisteredClaims 填充默认的 RegisteredClaims 字段
func PrepareRegisteredClaims(rc *jwt.RegisteredClaims) {
	if rc == nil {
		return
	}
	now := time.Now()
	if rc.Issuer == "" {
		rc.Issuer = defaultIssuer
	}
	if rc.Subject == "" {
		rc.Subject = defaultSubject
	}
	if len(rc.Audience) == 0 {
		rc.Audience = jwt.ClaimStrings{defaultAudience}
	}
	if rc.ID == "" {
		rc.ID = id.GenerateID()
	}
	if rc.IssuedAt == nil {
		rc.IssuedAt = jwt.NewNumericDate(now)
	}
	if rc.NotBefore == nil {
		rc.NotBefore = jwt.NewNumericDate(now)
	}
	if rc.ExpiresAt == nil && config.JWTExpiration > 0 {
		rc.ExpiresAt = jwt.NewNumericDate(now.Add(config.JWTExpiration))
	}
}

// SignHS256 使用 HS256 对 claims 进行签名
func SignHS256(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTKey))
}

// ParseHS256 使用 HS256 验证并解析 token，结果写入传入的 claims
func ParseHS256(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTKey), nil
	})
}

// JWTIssue 签发 JWT Token，使用 map 存储数据
// 参数 data 可以是任何 map[string]interface{}，完全灵活
// 使用示例：
//
//	data := map[string]interface{}{
//	    "user_id": "123",
//	    "username": "john",
//	    "role": "admin",
//	    "permissions": []string{"read", "write"},
//	}
//	token, err := JWTIssue(data)
func JWTIssue(data map[string]interface{}) (string, error) {
	claims := MapClaims{Data: data}
	PrepareRegisteredClaims(&claims.RegisteredClaims)
	return SignHS256(&claims)
}

// JWTDecrypt 解析 JWT Token，返回 map 数据
func JWTDecrypt(tokenString string) (map[string]interface{}, error) {
	claims := &MapClaims{}
	token, err := ParseHS256(tokenString, claims)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims.Data, nil
}
