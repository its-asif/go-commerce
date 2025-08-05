package models

import "time"

type CartItem struct {
	// User ID
	UserID int `db:"user_id" json:"user_id" example:"1" description:"ID of the user who owns the cart"`
	// Product ID
	ProductID int `db:"product_id" json:"product_id" example:"101" description:"ID of the product added to the cart"`
	// Quantity of the product
	Quantity int `db:"quantity" json:"quantity" example:"2" description:"Quantity of the product in the cart"`
	// Price of the product at the time it was added
	Price float64 `db:"price" json:"price" example:"199.99" description:"Price of the product at the time it was added"`
	// Timestamp when the product was added to the cart
	AddedAt time.Time `db:"added_at" json:"added_at" example:"2023-08-01T12:00:00Z" description:"Timestamp when the product was added to the cart"`
}

// CreateCartItemRequest represents the input for adding a product to the cart.
// @Description Input for adding a product to the cart.
type CreateCartItemRequest struct {
	// Product ID
	ProductID int `json:"product_id" example:"101" description:"ID of the product to add to the cart"`
	// Quantity of the product
	Quantity int `json:"quantity" example:"2" description:"Quantity of the product to add to the cart"`
}

// UpdateCartItemRequest represents the input for updating a product in the cart.
// @Description Input for updating the quantity of a product in the cart.
type UpdateCartItemRequest struct {
	// Quantity of the product
	Quantity int `json:"quantity" example:"3" description:"Updated quantity of the product in the cart"`
}

// CartItemResponse represents the details of a cart item in the response.
// @Description Details of a cart item in the response.
type CartItemResponse struct {
	// User ID
	UserID int `json:"user_id" example:"1" description:"ID of the user who owns the cart"`
	// Product ID
	ProductID int `json:"product_id" example:"101" description:"ID of the product in the cart"`
	// Quantity of the product
	Quantity int `json:"quantity" example:"2" description:"Quantity of the product in the cart"`
	// Price of the product at the time it was added
	Price float64 `json:"price" example:"199.99" description:"Price of the product at the time it was added"`
	// Timestamp when the product was added to the cart
	AddedAt time.Time `json:"added_at" example:"2023-08-01T12:00:00Z" description:"Timestamp when the product was added to the cart"`
}
