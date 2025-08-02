package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input models.Category
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO categories (name, slug)
				VALUES($1,$2)
				RETURNING id`
	err = db.DB.QueryRowx(query, input.Name, input.Slug).Scan(&input.ID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Invalidate categories cache
	_ = utils.DeleteCache("all_categories")

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(input)
}

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category

	// Check cache first
	cacheKey := "all_categories"
	err := utils.GetCache(cacheKey, &categories)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(categories)
		return
	}

	// Get from database if not in cache
	query := `SELECT * FROM categories`
	err = db.DB.Select(&categories, query)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Cache the categories
	_ = utils.SetCache(cacheKey, categories, time.Minute*30)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}
