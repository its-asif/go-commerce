package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/middleware"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

func TestGetCarts_CacheHit(t *testing.T) {
	orig := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error {
		dest := destination.(*[]models.CartItem)
		*dest = []models.CartItem{{ProductID: 2, Quantity: 1, Price: 3.5}}
		return nil
	}
	defer func() { utils.GetCache = orig }()

	req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
	ctxReq := req.WithContext(context.WithValue(req.Context(), middleware.UserKey, 10))
	rr := httptest.NewRecorder()
	GetCarts(rr, ctxReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var items []models.CartItem
	if err := json.NewDecoder(rr.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) != 1 || items[0].ProductID != 2 {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestAddToCart_Success(t *testing.T) {
	origGetProd := db.GetSingleProduct
	db.GetSingleProduct = func(id int) (models.Product, error) { return models.Product{ID: id, Price: 4.5}, nil }
	defer func() { db.GetSingleProduct = origGetProd }()

	origAdd := db.AddOrUpdateCartItem
	db.AddOrUpdateCartItem = func(userID, productID, quantity int, price float64) error { return nil }
	defer func() { db.AddOrUpdateCartItem = origAdd }()

	origDelCache := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDelCache }()

	body := []byte(`{"ProductID":2,"Quantity":1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/cart", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, 20))
	rr := httptest.NewRecorder()
	AddToCart(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}

func TestRemoveFromCart_Success(t *testing.T) {
	origDel := db.DeleteCartItem
	db.DeleteCartItem = func(userID, productID int) error { return nil }
	defer func() { db.DeleteCartItem = origDel }()

	origDelCache := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDelCache }()

	req := httptest.NewRequest(http.MethodDelete, "/api/cart/2", nil)
	req = mux.SetURLVars(req, map[string]string{"product_id": "2"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, 20))
	rr := httptest.NewRecorder()
	RemoveFromCart(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}
