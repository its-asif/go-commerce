package models

import "time"

// Order represents an order placed by a user.
// @Description An order contains details like user ID, total price, status, and the time it was placed.
type Order struct {
	// Order ID
	ID int `db:"id" json:"id" example:"1" description:"Unique identifier for the order"`
	// User ID
	UserID int `db:"user_id" json:"user_id" example:"101" description:"ID of the user who placed the order"`
	// Total Price of the order
	TotalPrice float64 `db:"total_price" json:"total_price" example:"299.99" description:"Total price of the order"`
	// Status of the order
	Status string `db:"status" json:"status" example:"Pending" description:"Status of the order (e.g., Pending, Completed, Cancelled)"`
	// Timestamp when the order was placed
	PlacedAt time.Time `db:"placed_at" json:"placed_at" example:"2023-08-01T12:00:00Z" description:"Timestamp when the order was placed"`
}

// OrderedItem represents an item in an order.
// @Description An ordered item contains details about the product, quantity, price, and the order it belongs to.
type OrderedItem struct {
	// Ordered Item ID
	ID int `db:"id" json:"id" example:"1" description:"Unique identifier for the ordered item"`
	// Order ID
	OrderID int `db:"order_id" json:"order_id" example:"1" description:"ID of the order this item belongs to"`
	// Product ID
	ProductID int `db:"product_id" json:"product_id" example:"101" description:"ID of the product in the order"`
	// Quantity of the product
	Quantity int `db:"quantity" json:"quantity" example:"2" description:"Quantity of the product in the order"`
	// Price of the product at the time of the order
	Price float64 `db:"price" json:"price" example:"149.99" description:"Price of the product at the time of the order"`
}

// OrderResponse represents the details of an order in the response.
// @Description Details of an order in the response.
type OrderResponse struct {
	// Order ID
	ID int `json:"id" example:"1" description:"Unique identifier for the order"`
	// User ID
	UserID int `json:"user_id" example:"101" description:"ID of the user who placed the order"`
	// Total Price of the order
	TotalPrice float64 `json:"total_price" example:"299.99" description:"Total price of the order"`
	// Status of the order
	Status string `json:"status" example:"Pending" description:"Status of the order"`
	// Timestamp when the order was placed
	PlacedAt time.Time `json:"placed_at" example:"2023-08-01T12:00:00Z" description:"Timestamp when the order was placed"`
	// List of ordered items
	Items []OrderedItem `json:"items" description:"List of items in the order"`
}
