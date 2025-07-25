package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/middleware"
	"github.com/its-asif/go-commerce/models"
	"net/http"
)

func Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)

	var cartItems []models.CartItem
	// get cart-items
	err := db.DB.Select(&cartItems, `SELECT * FROM cart_items WHERE user_id = $1`, userID)
	if err != nil || len(cartItems) == 0 {
		http.Error(w, "Cart is empty or unavailable", http.StatusBadRequest)
		return
	}

	// calculate total price
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}

	// setting order data
	order := models.Order{
		UserID:     userID,
		TotalPrice: total,
		Status:     "pending",
	}
	// Insert into orders
	err = db.DB.QueryRowx(
		`INSERT INTO orders (user_id, total_price, status)
		 VALUES ($1, $2, $3)
		 RETURNING id, placed_at`,
		order.UserID, order.TotalPrice, order.Status,
	).Scan(&order.ID, &order.PlacedAt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	// Insert ordered items
	for _, item := range cartItems {
		_, err := db.DB.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to save order items", http.StatusInternalServerError)
			return
		}
	}

	// Clear cart
	_, _ = db.DB.Exec(`DELETE FROM cart_items WHERE user_id = $1`, userID)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)

	var orders []models.Order
	err := db.DB.Select(&orders, `SELECT * FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
