package models

import "time"

type Order struct {
	ID       int       `db:"id" json:"id"`
	UserID   int       `db:"user_id" json:"user_id"`
	Total    float64   `db:"total" json:"total"`
	Status   string    `db:"status" json:"status"`
	PlacedAt time.Time `db:"placed_at" json:"placed_at"`
}

type OrderedItem struct {
	ID        int     `db:"id" json:"id"`
	OrderID   int     `db:"order_id" json:"order_id"`
	ProductID int     `db:"product_id" json:"product_id"`
	Quantity  int     `db:"qty" json:"qty"`
	Price     float64 `db:"price" json:"price"`
}
