package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/its-asif/go-commerce/utils"
)

func TestAuthMiddlewareValid(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	tok, err := utils.GenerateTokenWithRole(7, "user")
	if err != nil {
		t.Fatalf("GenerateTokenWithRole error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok))
	rr := httptest.NewRecorder()

	next := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(UserKey).(int)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, uid)
	}))

	next.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	if rr.Body.String() != "7" {
		t.Fatalf("expected body '7' got %q", rr.Body.String())
	}
}

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
	var payload map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload["error"] == "" {
		t.Fatalf("expected error message in response")
	}
}

func TestAuthMiddlewareInvalidTokenAndMissingUid(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	// invalid token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	rr := httptest.NewRecorder()
	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token got %d", rr.Code)
	}

	// token missing uid claim: craft token that only contains a role
	tokenOnlyRole := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "user",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	})
	tkn, err := tokenOnlyRole.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("failed to sign tokenOnlyRole: %v", err)
	}
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("Authorization", "Bearer "+tkn)
	rr2 := httptest.NewRecorder()
	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for token missing uid got %d", rr2.Code)
	}
}
