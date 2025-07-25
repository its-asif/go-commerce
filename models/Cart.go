package models

import "time"

type CartItem struct {
	UserID    int       `db:"user_id" json:"user_id"`
	ProductID int       `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Price     float64   `db:"price" json:"price"`
	AddedAt   time.Time `db:"added_at" json:"added_at"`
}
