package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sgitwhyd/jagong/pkg/env"
	"log"
	"time"
)

type ClaimToken struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

var mapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 24 * 3,
}

func GenerateToken(username string, fullname string, tokenType string) (string, error) {
	secret := []byte(env.GetEnv("JWT_SECRET", ""))

	claimToken := ClaimToken{
		Username: username,
		FullName: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	resultToken, err := token.SignedString(secret)
	if err != nil {
		log.Printf("token.SignedWithClaims err:%v", err)
		return "", err
	}

	return resultToken, nil
}
