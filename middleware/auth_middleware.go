package middleware

import (
	"context"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/utils"
	"net/http"
	"strings"
)

const UserKey int = 0

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := utils.ParseToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, int(claims["uid"].(float64)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var role string
		err := db.DB.Get(&role, "SELECT role FROM users WHERE id=$1", userID)
		if err != nil || role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
