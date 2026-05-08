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

// @Summary      Get Cart Items
// @Description  Retrieve all items in the user's cart
// @Tags         Cart
// @Produce      json
// @Param        Authorization header string true "Bearer + JWT_Token>"
// @Success      200 {array} models.CartItemResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Server Error"
// @Router       /api/cart [get]
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

// @Summary      Add to Cart
// @Description  Add a product to the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer + JWT_Token"
// @Param        input body models.CreateCartItemRequest true "Cart Item Input"
// @Success      201 {object} models.CartItemResponse
// @Failure      400 {string} string "Invalid Input"
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Server Error"
// @Router       /api/cart [post]
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
	// attempt to read price via DB directly
	// reuse GetSingleProduct from db helpers to obtain product price
	prod, pErr := db.GetSingleProduct(input.ProductID)
	if pErr != nil {
		http.Error(w, "Product not found", http.StatusBadRequest)
		return
	}
	price = prod.Price

	if err := db.AddOrUpdateCartItem(userID, input.ProductID, input.Quantity, price); err != nil {
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

// @Summary      Remove from cart
// @Description  Remove a product to the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer + JWT_Token"
// @Success      201 {object} models.CartItemResponse
// @Failure      400 {string} string "Invalid Input"
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Server Error"
// @Router       /api/cart [post]
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(int)
	prodID := mux.Vars(r)["product_id"]

	productID, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteCartItem(userID, productID); err != nil {
		http.Error(w, "Failed to remove item", http.StatusInternalServerError)
		return
	}

	// Invalidate cart cache
	cacheKey := fmt.Sprintf("cart_user_%d", userID)
	_ = utils.DeleteCache(cacheKey)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Product removed from cart")
}
