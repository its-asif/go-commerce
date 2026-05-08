package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

func TestGetAllProducts_DBReturned(t *testing.T) {
	// stub cache to simulate miss
	origGetCache := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error { return fmtError("cache miss") }
	defer func() { utils.GetCache = origGetCache }()
	origSetCache := utils.SetCache
	utils.SetCache = func(key string, value interface{}, expiration time.Duration) error { return nil }
	defer func() { utils.SetCache = origSetCache }()

	// stub DB
	origGetAll := db.GetAllProduct
	db.GetAllProduct = func() ([]models.Product, error) {
		return []models.Product{{ID: 1, Name: "P1", Price: 9.99}}, nil
	}
	defer func() { db.GetAllProduct = origGetAll }()

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rr := httptest.NewRecorder()
	GetAllProducts(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var prods []models.Product
	if err := json.NewDecoder(rr.Body).Decode(&prods); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(prods) != 1 || prods[0].ID != 1 {
		t.Fatalf("unexpected products: %#v", prods)
	}
}

func TestGetOneProduct_DBReturned(t *testing.T) {
	origGetCache := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error { return fmtError("cache miss") }
	defer func() { utils.GetCache = origGetCache }()
	origSetCache := utils.SetCache
	utils.SetCache = func(key string, value interface{}, expiration time.Duration) error { return nil }
	defer func() { utils.SetCache = origSetCache }()

	origGetSingle := db.GetSingleProduct
	db.GetSingleProduct = func(id int) (models.Product, error) {
		return models.Product{ID: id, Name: "P1"}, nil
	}
	defer func() { db.GetSingleProduct = origGetSingle }()

	bodyReq := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	bodyReq = mux.SetURLVars(bodyReq, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	GetOneProduct(rr, bodyReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var p models.Product
	if err := json.NewDecoder(rr.Body).Decode(&p); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if p.ID != 1 {
		t.Fatalf("unexpected product id: %d", p.ID)
	}
}

func TestCreateProducts_Success(t *testing.T) {
	origCreate := db.CreateProduct
	db.CreateProduct = func(input models.Product) error { return nil }
	defer func() { db.CreateProduct = origCreate }()

	origDel := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDel }()

	p := models.Product{Name: "New", Price: 5.0}
	b, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	CreateProducts(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}

func TestUpdateOneProduct_Success(t *testing.T) {
	origUpdate := db.UpdateProduct
	db.UpdateProduct = func(id int, input models.UpdateProductRequest, w http.ResponseWriter) error { return nil }
	defer func() { db.UpdateProduct = origUpdate }()

	origDel := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDel }()

	name := "Updated"
	upd := models.UpdateProductRequest{Name: &name}
	b, _ := json.Marshal(upd)
	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	UpdateOneProduct(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestDeleteOneProduct_Success(t *testing.T) {
	origDelete := db.DeleteProduct
	db.DeleteProduct = func(id int, w http.ResponseWriter) error { return nil }
	defer func() { db.DeleteProduct = origDelete }()

	origDel := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDel }()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	DeleteOneProduct(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetAllProducts_CacheHit(t *testing.T) {
	origGetCache := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error {
		dest := destination.(*[]models.Product)
		*dest = []models.Product{{ID: 3, Name: "pc"}}
		return nil
	}
	defer func() { utils.GetCache = origGetCache }()

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rr := httptest.NewRecorder()
	GetAllProducts(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetOneProduct_CacheHit(t *testing.T) {
	origGetCache := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error {
		dest := destination.(*models.Product)
		*dest = models.Product{ID: 5, Name: "p5"}
		return nil
	}
	defer func() { utils.GetCache = origGetCache }()

	req := httptest.NewRequest(http.MethodGet, "/api/products/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rr := httptest.NewRecorder()
	GetOneProduct(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestCreateProducts_BadInput(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewReader([]byte("notjson")))
	rr := httptest.NewRecorder()
	CreateProducts(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestUpdateOneProduct_BadInput(t *testing.T) {
	origDelCache := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDelCache }()

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewReader([]byte("notjson")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	UpdateOneProduct(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestDeleteOneProduct_BadID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/products/notint", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "notint"})
	rr := httptest.NewRecorder()
	DeleteOneProduct(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

// helper to create an error implementing error without importing fmt in many places
type errString string

func (e errString) Error() string { return string(e) }

func fmtError(s string) error { return errString(s) }
