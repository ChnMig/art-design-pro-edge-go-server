package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"api-server/config"
	"api-server/util/id"
)

// MultiTenantClaims 多租户JWT Claims
type MultiTenantClaims struct {
	UserID   uint   `json:"user_id"`
	TenantID uint   `json:"tenant_id"`
	Account  string `json:"account"`
	jwt.RegisteredClaims
}

// JWTIssue 签发多租户JWT token
func JWTIssue(userID, tenantID uint, account string) (string, error) {
	// set key
	mySigningKey := []byte(config.JWTKey)
	// Calculate expiration time
	nt := time.Now()
	exp := nt.Add(config.JWTExpiration)
	// Create the Claims
	claims := MultiTenantClaims{
		UserID:   userID,
		TenantID: tenantID,
		Account:  account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "server",
			IssuedAt:  jwt.NewNumericDate(nt),
			Subject:   "token",
			Audience:  jwt.ClaimStrings{"client"},
			NotBefore: jwt.NewNumericDate(nt),
			ID:        id.IssueMd5ID(),
		},
	}
	// issue
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	st, err := t.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return st, nil
}

// JWTDecrypt 解析多租户JWT token
func JWTDecrypt(tokenString string) (*MultiTenantClaims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &MultiTenantClaims{}, func(token *jwt.Token) (interface{}, error) {
		// HMAC Check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims, ok := t.Claims.(*MultiTenantClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid claims")
	}
}

// EncodeUserInfo 将用户和租户信息编码为JSON字符串（用于兼容原有系统）
func EncodeUserInfo(userID, tenantID uint, account string) (string, error) {
	userInfo := map[string]interface{}{
		"user_id":   userID,
		"tenant_id": tenantID,
		"account":   account,
	}
	data, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DecodeUserInfo 解码用户信息（用于兼容原有系统）
func DecodeUserInfo(data string) (userID, tenantID uint, account string, err error) {
	var userInfo map[string]interface{}
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		return 0, 0, "", err
	}
	
	if uid, ok := userInfo["user_id"].(float64); ok {
		userID = uint(uid)
	}
	if tid, ok := userInfo["tenant_id"].(float64); ok {
		tenantID = uint(tid)
	}
	if acc, ok := userInfo["account"].(string); ok {
		account = acc
	}
	
	return userID, tenantID, account, nil
}