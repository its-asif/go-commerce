package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/its-asif/go-commerce/db"
)

func TestAdminMiddleware_AllowsAdmin(t *testing.T) {
	orig := db.GetUserRole
	db.GetUserRole = func(userID int) (string, error) { return "admin", nil }
	defer func() { db.GetUserRole = orig }()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserKey, 1))
	rr := httptest.NewRecorder()
	handler := AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestAdminMiddleware_Forbidden(t *testing.T) {
	orig := db.GetUserRole
	db.GetUserRole = func(userID int) (string, error) { return "user", nil }
	defer func() { db.GetUserRole = orig }()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserKey, 2))
	rr := httptest.NewRecorder()
	handler := AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rr.Code)
	}
}
