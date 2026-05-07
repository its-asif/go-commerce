package models

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Role      string    `db:"role" json:"role"` // user or admin
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Optional / new fields
	Address   *string    `db:"address" json:"address,omitempty"`
	Phone     *string    `db:"phone" json:"phone,omitempty"`
	Avatar    *string    `db:"avatar" json:"avatar,omitempty"`
	LastLogin *time.Time `db:"last_login" json:"last_login,omitempty"`
	Status    string     `db:"status" json:"status,omitempty"` // active, inactive, suspended
}

type UpdateProfileRequest struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
	Avatar  *string `json:"avatar"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
