package authentication

import (
	"fmt"
	"time"

	"api-server/config"
	"api-server/util/id"

	"github.com/golang-jwt/jwt/v5"
)

var issue string = "server"
var subject string = "token"
var audienc string = "client"

type MyCustomClaims struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

// JWTIssue issue jwt
func JWTIssue(data string) (string, error) {
	// set key
	mySigningKey := []byte(config.JWTKey)
	// Calculate expiration time
	nt := time.Now()
	exp := nt.Add(config.JWTExpiration)
	// Create the Claims
	claims := MyCustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    issue,
			IssuedAt:  jwt.NewNumericDate(nt),
			Subject:   subject,
			Audience:  jwt.ClaimStrings{audienc},
			NotBefore: jwt.NewNumericDate(nt),
			ID:        id.IssueMd5ID(),
		},
	}
	// https://en.wikipedia.org/wiki/JSON_Web_Token
	// issue
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	st, err := t.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return st, nil
}

// JWTDecrypt string token to data
func JWTDecrypt(tokenString string) (string, error) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// HMAC Check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTKey), nil
	})
	if err != nil {
		return "", err
	}
	if !t.Valid {
		return "", fmt.Errorf("invalid token")
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims["data"].(string), nil
	} else {
		return "", err
	}
}
