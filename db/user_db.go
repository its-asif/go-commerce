package db

import (
	"fmt"

	"github.com/its-asif/go-commerce/models"
)

var GetAllUsers = func() ([]models.User, error) {
	var user []models.User
	query := `SELECT * FROM users`
	err := DB.Select(&user, query)
	return user, err
}

var GetSingleUser = func(key string, value interface{}) (models.User, error) {
	var user models.User
	var query string
	switch key {
	case "id":
		query = `SELECT * FROM users WHERE id = $1`
	case "email":
		query = `SELECT * FROM users WHERE email = $1`
	default:
		return user, fmt.Errorf("invalid key: %s", key)
	}
	err := DB.Get(&user, query, value)
	return user, err
}

// GetUserRole fetches the role string for a user id. Overridable for tests.
var GetUserRole = func(userID int) (string, error) {
	var role string
	err := DB.Get(&role, "SELECT role FROM users WHERE id=$1", userID)
	return role, err
}
