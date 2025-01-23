package jwt_token

import (
	"context"
	"fmt"
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

var MapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 24 * 3,
}

var JWT_SECRET = []byte(env.GetEnv("JWT_SECRET", ""))

func GenerateToken(username string, fullname string, tokenType string) (string, error) {
	claimToken := ClaimToken{
		Username: username,
		FullName: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(MapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	resultToken, err := token.SignedString(JWT_SECRET)
	if err != nil {
		log.Printf("jwt_token.SignedWithClaims err:%v", err)
		return "", err
	}

	return resultToken, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		log.Printf("failed to parse token: %v", err)
		return nil, err
	}

	// Validasi jwtToken dan claimToken
	if jwtToken == nil || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claimToken, ok := jwtToken.Claims.(*ClaimToken)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claimToken, nil
}
