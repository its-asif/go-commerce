package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/middleware"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

func GetCarts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)

	var cartItems []models.CartItem

	// Check cache first
	cacheKey := fmt.Sprintf("cart_user_%d", userID)
	err := utils.GetCache(cacheKey, &cartItems)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cartItems)
		return
	}

	// Get from database if not in cache
	query := `SELECT * FROM cart_items WHERE user_id = $1`
	err = db.DB.Select(&cartItems, query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
		return
	}

	// Cache the cart items
	_ = utils.SetCache(cacheKey, cartItems, time.Minute*5)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cartItems)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)

	var input models.CartItem
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get price from product table
	var price float64
	err = db.DB.Get(&price, `SELECT price FROM products WHERE id = $1`, input.ProductID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO cart_items (user_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`
	_, err = db.DB.Exec(query, userID, input.ProductID, input.Quantity, price)
	if err != nil {
		log.Println("DB Error:", err)
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	// Invalidate cart cache
	cacheKey := fmt.Sprintf("cart_user_%d", userID)
	_ = utils.DeleteCache(cacheKey)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Product added to cart")
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)
	prodID := mux.Vars(r)["product_id"]

	productID, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`
	_, err = db.DB.Exec(query, userID, productID)
	if err != nil {
		http.Error(w, "Failed to remove item", http.StatusInternalServerError)
		return
	}

	// Invalidate cart cache
	cacheKey := fmt.Sprintf("cart_user_%d", userID)
	_ = utils.DeleteCache(cacheKey)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Product removed from cart")
}
