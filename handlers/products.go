package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"net/http"
	"strconv"
)

func CreateProducts(w http.ResponseWriter, r *http.Request) {
	var input models.Product
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (name, description, price, stock, category_id, image_url)
				VALUES($1,$2,$3,$4,$5,$6)
				RETURNING id, created_at`
	err = db.DB.QueryRowx(query, input.Name, input.Description, input.Price, input.Stock, input.CategoryID, input.ImageURL).Scan(&input.ID, &input.CreatedAt)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(input)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var prod []models.Product
	query := `SELECT * FROM products`
	err := db.DB.Select(&prod, query)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if len(prod) == 0 {
		json.NewEncoder(w).Encode("No tasks found")
		return
	}

	w.WriteHeader(http.StatusFound)
	_ = json.NewEncoder(w).Encode(prod)
}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var prod models.Product
	query := `SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`
	row := db.DB.QueryRow(query, id)
	err = row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CategoryID, &prod.ImageURL)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusFound)
	_ = json.NewEncoder(w).Encode(prod)
}

func UpdateOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var prod models.Product
	err = json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}

	//var result models.Product
	query := `UPDATE products 
				SET name=$1, description=$2, price=$3, stock=$4, category_id=$5, image_url=$6
				WHERE id=$7`

	_, err = db.DB.Exec(query, prod.Name, prod.Description, prod.Price, prod.Stock, prod.CategoryID, prod.ImageURL, id)
	if err != nil {
		http.Error(w, "Error updating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(prod)
}

func DeleteOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	query := `DELETE FROM products 
				WHERE id=$1`

	_, err = db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error updating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode("Deleted Successfully")
}
