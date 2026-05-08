package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

func TestGetAllCategories_CacheHit(t *testing.T) {
	orig := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error {
		dest := destination.(*[]models.Category)
		*dest = []models.Category{{ID: 1, Name: "c"}}
		return nil
	}
	defer func() { utils.GetCache = orig }()

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	rr := httptest.NewRecorder()
	GetAllCategories(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var cats []models.Category
	if err := json.NewDecoder(rr.Body).Decode(&cats); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(cats) != 1 || cats[0].ID != 1 {
		t.Fatalf("unexpected cats: %#v", cats)
	}
}

func TestCreateCategory_Success(t *testing.T) {
	orig := db.CreateCategory
	db.CreateCategory = func(input models.Category) (models.Category, error) {
		input.ID = 77
		return input, nil
	}
	defer func() { db.CreateCategory = orig }()

	origDelCache := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDelCache }()

	body := []byte(`{"Name":"c","Slug":"c"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	CreateCategory(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}
