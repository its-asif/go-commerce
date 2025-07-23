package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/its-asif/go-commerce/config"
	"time"
)

func GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	return tokenString, err
}

func ParseToken(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})

	if err != nil || !t.Valid {
		return nil, err
	}

	return t.Claims.(jwt.MapClaims), err
}
