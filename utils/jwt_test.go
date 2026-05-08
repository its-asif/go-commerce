package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	tok, err := GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	claims, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}

	uidAny, ok := claims["uid"]
	if !ok {
		t.Fatalf("uid claim missing")
	}
	if v, ok := uidAny.(float64); !ok || int(v) != 42 {
		t.Fatalf("unexpected uid claim: %#v", uidAny)
	}
}

func TestGenerateTokenWithRoleAndParse(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	tok, err := GenerateTokenWithRole(7, "admin")
	if err != nil {
		t.Fatalf("GenerateTokenWithRole error: %v", err)
	}

	claims, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}

	if role, ok := claims["role"].(string); !ok || role != "admin" {
		t.Fatalf("unexpected role claim: %#v", claims["role"])
	}
}

func TestParseExpiredAndMalformedToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	// expired token
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": 1,
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	})
	s, err := expired.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}
	if _, err := ParseToken(s); err == nil {
		t.Fatalf("expected error for expired token")
	}

	// malformed token
	if _, err := ParseToken("not-a-token"); err == nil {
		t.Fatalf("expected error for malformed token")
	}
}
