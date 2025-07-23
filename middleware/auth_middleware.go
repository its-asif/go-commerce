package middleware

import (
	"context"
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
