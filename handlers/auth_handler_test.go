package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
	"github.com/lib/pq"
)

func TestLogin_CacheHit(t *testing.T) {
	// stub GetCache to return user in cache
	origGetCache := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error {
		u := destination.(*models.User)
		*u = models.User{ID: 5, Email: "a@b.c", Password: "hashed", Role: "user"}
		return nil
	}
	defer func() { utils.GetCache = origGetCache }()

	origMatch := utils.MatchPass
	utils.MatchPass = func(hashedPass, pass string) error { return nil }
	defer func() { utils.MatchPass = origMatch }()

	origGen := utils.GenerateTokenWithRole
	utils.GenerateTokenWithRole = func(userID int, role string) (string, error) { return "tkn", nil }
	defer func() { utils.GenerateTokenWithRole = origGen }()

	body := []byte(`{"Email":"a@b.c","Password":"pw"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	Login(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["token"] != "tkn" {
		t.Fatalf("unexpected token: %v", resp)
	}
}

func TestRegister_Success(t *testing.T) {
	origCreate := db.CreateUser
	db.CreateUser = func(name, email, hashedPassword string) (models.User, error) {
		return models.User{ID: 11, Email: email}, nil
	}
	defer func() { db.CreateUser = origCreate }()

	body := []byte(`{"Name":"A","Email":"x@x","Password":"pw"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	Register(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}

func TestLogin_DBHit(t *testing.T) {
	// simulate cache miss
	origGetCache := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error { return errString("miss") }
	defer func() { utils.GetCache = origGetCache }()

	origSetCache := utils.SetCache
	utils.SetCache = func(key string, value interface{}, expiration time.Duration) error { return nil }
	defer func() { utils.SetCache = origSetCache }()

	origGetUser := db.GetUserByEmail
	db.GetUserByEmail = func(email string) (models.User, error) {
		return models.User{ID: 6, Email: email, Password: "h"}, nil
	}
	defer func() { db.GetUserByEmail = origGetUser }()

	origMatch := utils.MatchPass
	utils.MatchPass = func(hashedPass, pass string) error { return nil }
	defer func() { utils.MatchPass = origMatch }()

	origGen := utils.GenerateTokenWithRole
	utils.GenerateTokenWithRole = func(userID int, role string) (string, error) { return "tkn2", nil }
	defer func() { utils.GenerateTokenWithRole = origGen }()

	body := []byte(`{"Email":"b@b","Password":"pw"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	Login(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("notjson")))
	rr := httptest.NewRecorder()
	Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestRegister_HashFail(t *testing.T) {
	origHash := utils.HashPass
	utils.HashPass = func(pw string) (string, error) { return "", errString("fail") }
	defer func() { utils.HashPass = origHash }()

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(`{"Name":"A","Email":"x@x","Password":"pw"}`)))
	rr := httptest.NewRecorder()
	Register(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	origCreate := db.CreateUser
	db.CreateUser = func(name, email, hashedPassword string) (models.User, error) {
		return models.User{}, &pq.Error{Code: "23505"}
	}
	defer func() { db.CreateUser = origCreate }()

	// ensure HashPass returns ok
	origHash := utils.HashPass
	utils.HashPass = func(pw string) (string, error) { return "h", nil }
	defer func() { utils.HashPass = origHash }()

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(`{"Name":"A","Email":"x@x","Password":"pw"}`)))
	rr := httptest.NewRecorder()
	Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}
