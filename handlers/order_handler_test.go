package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/middleware"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

func TestGetOrders_CacheHit(t *testing.T) {
	orig := utils.GetCache
	utils.GetCache = func(key string, destination interface{}) error {
		dest := destination.(*[]models.Order)
		*dest = []models.Order{{ID: 1, TotalPrice: 10.0}}
		return nil
	}
	defer func() { utils.GetCache = orig }()

	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	ctxReq := req.WithContext(context.WithValue(req.Context(), middleware.UserKey, 11))
	rr := httptest.NewRecorder()
	GetOrders(rr, ctxReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var orders []models.Order
	if err := json.NewDecoder(rr.Body).Decode(&orders); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(orders) != 1 || orders[0].ID != 1 {
		t.Fatalf("unexpected orders: %#v", orders)
	}
}

func TestCheckout_Success(t *testing.T) {
	origGetCart := db.GetCartItems
	db.GetCartItems = func(userID int) ([]models.CartItem, error) {
		return []models.CartItem{{ProductID: 2, Quantity: 1, Price: 5.0}}, nil
	}
	defer func() { db.GetCartItems = origGetCart }()

	origInsertOrder := db.InsertOrder
	db.InsertOrder = func(order *models.Order) error {
		order.ID = 55
		return nil
	}
	defer func() { db.InsertOrder = origInsertOrder }()

	origInsertItem := db.InsertOrderItem
	db.InsertOrderItem = func(orderID, productID, quantity int, price float64) error { return nil }
	defer func() { db.InsertOrderItem = origInsertItem }()

	origDelete := db.DeleteCartByUser
	db.DeleteCartByUser = func(userID int) error { return nil }
	defer func() { db.DeleteCartByUser = origDelete }()

	origDelCache := utils.DeleteCache
	utils.DeleteCache = func(key string) error { return nil }
	defer func() { utils.DeleteCache = origDelCache }()

	req := httptest.NewRequest(http.MethodPost, "/api/orders/checkout", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, 33))
	rr := httptest.NewRecorder()
	Checkout(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}
