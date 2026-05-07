package db

import (
	"fmt"
	"github.com/its-asif/go-commerce/models"
)

func GetAllUsers() ([]models.User, error) {
	var user []models.User
	query := `SELECT * FROM users`
	err := DB.Select(&user, query)
	return user, err
}

func GetSingleUser(key string, value interface{}) (models.User, error) {
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
