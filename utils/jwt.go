package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/its-asif/go-commerce/config"
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

// GenerateTokenWithRole issues a JWT that includes user ID and role.
func GenerateTokenWithRole(userID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})
	return token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
}

// ParseToken validates the JWT, enforces HS256 algorithm, and checks expiration.
func ParseToken(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// Enforce expected signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Validate exp if present (jwt/v5 validates with Parse + RegisteredClaims, but we use MapClaims)
	if expRaw, ok := claims["exp"]; ok {
		switch v := expRaw.(type) {
		case float64:
			if time.Unix(int64(v), 0).Before(time.Now()) {
				return nil, errors.New("token expired")
			}
		case int64:
			if time.Unix(v, 0).Before(time.Now()) {
				return nil, errors.New("token expired")
			}
		}
	}
	return claims, nil
}
