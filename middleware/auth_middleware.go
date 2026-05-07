package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/utils"
)

// contextKey prevents collisions with other context keys
type contextKey string

const (
	UserKey contextKey = "userID"
	RoleKey contextKey = "role"
)

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			writeJSONError(w, http.StatusUnauthorized, "missing token")
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := utils.ParseToken(token)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Extract user ID safely
		uidAny, ok := claims["uid"]
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "missing uid claim")
			return
		}
		var userID int
		switch v := uidAny.(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		default:
			writeJSONError(w, http.StatusUnauthorized, "invalid uid claim type")
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, userID)
		if roleAny, ok := claims["role"]; ok {
			if roleStr, ok := roleAny.(string); ok {
				ctx = context.WithValue(ctx, RoleKey, roleStr)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prefer role from context (if token carries it), else fallback to DB check.
		if role, ok := r.Context().Value(RoleKey).(string); ok && role == "admin" {
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := r.Context().Value(UserKey).(int)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var role string
		err := db.DB.Get(&role, "SELECT role FROM users WHERE id=$1", userID)
		if err != nil || role != "admin" {
			writeJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}
