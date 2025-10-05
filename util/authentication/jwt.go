package authentication

import (
	"fmt"
	"time"

	"api-server/config"
	"api-server/util/id"

	"github.com/golang-jwt/jwt/v5"
)

var (
	defaultIssuer   = "server"
	defaultSubject  = "token"
	defaultAudience = "client"
)

type MyCustomClaims struct {
	Data string `json:"data"`
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
		rc.ID = id.IssueMd5ID()
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

// JWTIssue issue jwt
func JWTIssue(data string) (string, error) {
	claims := MyCustomClaims{Data: data}
	PrepareRegisteredClaims(&claims.RegisteredClaims)
	return SignHS256(&claims)
}

// JWTDecrypt string token to data
func JWTDecrypt(tokenString string) (string, error) {
	claims := &MyCustomClaims{}
	token, err := ParseHS256(tokenString, claims)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.Data, nil
}
