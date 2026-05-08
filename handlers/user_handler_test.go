package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
)

func TestGetAllUsers_Stubbed(t *testing.T) {
	orig := db.GetAllUsers
	db.GetAllUsers = func() ([]models.User, error) {
		return []models.User{{ID: 2, Email: "x@y"}}, nil
	}
	defer func() { db.GetAllUsers = orig }()

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rr := httptest.NewRecorder()
	GetAllUsers(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var users []models.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(users) != 1 || users[0].ID != 2 {
		t.Fatalf("unexpected users: %#v", users)
	}
}

func TestGetSingleUserByID_Stubbed(t *testing.T) {
	orig := db.GetSingleUser
	db.GetSingleUser = func(key string, value interface{}) (models.User, error) {
		return models.User{ID: value.(int), Email: "a@b"}, nil
	}
	defer func() { db.GetSingleUser = orig }()

	req := httptest.NewRequest(http.MethodGet, "/api/users/id/3", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "3"})
	rr := httptest.NewRecorder()
	GetSingleUserByID(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var u models.User
	if err := json.NewDecoder(rr.Body).Decode(&u); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if u.ID != 3 {
		t.Fatalf("unexpected user id: %d", u.ID)
	}
}

func TestGetSingleUserByEmail_Stubbed(t *testing.T) {
	orig := db.GetSingleUser
	db.GetSingleUser = func(key string, value interface{}) (models.User, error) {
		return models.User{ID: 9, Email: value.(string)}, nil
	}
	defer func() { db.GetSingleUser = orig }()

	req := httptest.NewRequest(http.MethodGet, "/api/users/email/test@x", nil)
	req = mux.SetURLVars(req, map[string]string{"email": "test@x"})
	rr := httptest.NewRecorder()
	GetSingleUserByEmail(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var u models.User
	if err := json.NewDecoder(rr.Body).Decode(&u); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if u.Email != "test@x" {
		t.Fatalf("unexpected email: %s", u.Email)
	}
}
