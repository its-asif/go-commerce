package db

import (
	"github.com/its-asif/go-commerce/models"
)

// CreateUser inserts a new user and returns the created user.
var CreateUser = func(name, email, hashedPassword string) (models.User, error) {
	user := models.User{}
	query := `INSERT INTO users(name, email, password) Values ($1, $2, $3) RETURNING id, created_at`
	err := DB.QueryRowx(query, name, email, hashedPassword).Scan(&user.ID, &user.CreatedAt)
	return user, err
}

// GetUserByEmail fetches a user by email.
var GetUserByEmail = func(email string) (models.User, error) {
	var user models.User
	err := DB.Get(&user, "Select * FROM users where email=$1", email)
	return user, err
}

// Category helpers
var GetAllCategories = func() ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT * FROM categories`
	err := DB.Select(&categories, query)
	return categories, err
}

var CreateCategory = func(input models.Category) (models.Category, error) {
	query := `INSERT INTO categories (name, slug) VALUES($1,$2) RETURNING id`
	err := DB.QueryRowx(query, input.Name, input.Slug).Scan(&input.ID)
	return input, err
}

// Cart helpers
var GetCartItems = func(userID int) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	query := `SELECT * FROM cart_items WHERE user_id = $1`
	err := DB.Select(&cartItems, query, userID)
	return cartItems, err
}

var AddOrUpdateCartItem = func(userID, productID, quantity int, price float64) error {
	query := `INSERT INTO cart_items (user_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`
	_, err := DB.Exec(query, userID, productID, quantity, price)
	return err
}

var DeleteCartItem = func(userID, productID int) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`
	_, err := DB.Exec(query, userID, productID)
	return err
}

// Order helpers
var GetOrdersByUser = func(userID int) ([]models.Order, error) {
	var orders []models.Order
	err := DB.Select(&orders, `SELECT * FROM orders WHERE user_id = $1`, userID)
	return orders, err
}

var InsertOrder = func(order *models.Order) error {
	return DB.QueryRowx(`INSERT INTO orders (user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id, placed_at`, order.UserID, order.TotalPrice, order.Status).Scan(&order.ID, &order.PlacedAt)
}

var InsertOrderItem = func(orderID, productID, quantity int, price float64) error {
	_, err := DB.Exec(`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`, orderID, productID, quantity, price)
	return err
}

var DeleteCartByUser = func(userID int) error {
	_, err := DB.Exec(`DELETE FROM cart_items WHERE user_id = $1`, userID)
	return err
}
