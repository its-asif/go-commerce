package models

import "time"

type CartItem struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	ProductID int       `db:"product_id" json:"product_id"`
	Quantity  int       `db:"qty" json:"qty"`
	Price     float64   `db:"price" json:"price"`
	AddedAt   time.Time `db:"added_at" json:"added_at"`
}
