package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/middleware"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

// @Summary      Place Order
// @Description  Place a new order with the items in the user's cart
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer + JWT_Token>"
// @Param        input body models.CreateOrderRequest true "Order Input"
// @Success      201 {object} models.OrderResponse
// @Failure      400 {string} string "Invalid Input"
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Server Error"
// @Router       /api/orders/checkout [post]
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

	// Invalidate user orders cache and cart cache
	ordersCacheKey := fmt.Sprintf("orders_user_%d", userID)
	cartCacheKey := fmt.Sprintf("cart_user_%d", userID)
	_ = utils.DeleteCache(ordersCacheKey)
	_ = utils.DeleteCache(cartCacheKey)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

// @Summary      Get Orders
// @Description  Retrieve all orders placed by the user
// @Tags         Orders
// @Produce      json
// @Param        Authorization header string true "Bearer JWT_Token"
// @Success      200 {array} models.OrderResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Server Error"
// @Router       /api/orders [get]
func GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)

	var orders []models.Order

	// Check cache first
	cacheKey := fmt.Sprintf("orders_user_%d", userID)
	err := utils.GetCache(cacheKey, &orders)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orders)
		return
	}

	// Get from database if not in cache
	err = db.DB.Select(&orders, `SELECT * FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}

	// Cache the orders
	_ = utils.SetCache(cacheKey, orders, time.Minute*10)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
