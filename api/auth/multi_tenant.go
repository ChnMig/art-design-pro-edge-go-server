package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"api-server/config"
	authutil "api-server/util/authentication"
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
	nt := time.Now()
	exp := nt.Add(config.JWTExpiration)
	claims := MultiTenantClaims{
		UserID:   userID,
		TenantID: tenantID,
		Account:  account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(nt),
			NotBefore: jwt.NewNumericDate(nt),
		},
	}
	authutil.PrepareRegisteredClaims(&claims.RegisteredClaims)
	return authutil.SignHS256(&claims)
}

// JWTDecrypt 解析多租户JWT token
func JWTDecrypt(tokenString string) (*MultiTenantClaims, error) {
	claims := &MultiTenantClaims{}
	t, err := authutil.ParseHS256(tokenString, claims)
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
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
