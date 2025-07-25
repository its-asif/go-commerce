package handlers

import (
	"encoding/json"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"net/http"
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

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(input)
}
