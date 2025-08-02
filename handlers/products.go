package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
	"net/http"
	"strconv"
	"time"
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
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	// invalid redis cache
	_ = utils.DeleteCache("all_products")

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(input)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var prod []models.Product

	// check in cache
	cacheKey := "all_products"
	err := utils.GetCache(cacheKey, &prod)
	if err == nil { // if no error found
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(prod)
		return
	}

	// get from DB
	query := `SELECT * FROM products`
	err = db.DB.Select(&prod, query)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	if len(prod) == 0 {
		json.NewEncoder(w).Encode("No tasks found")
		return
	}

	// setting cache
	_ = utils.SetCache(cacheKey, prod, time.Minute*10)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(prod)
}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	var prod models.Product

	// check if cache
	cacheId := fmt.Sprintf("product_" + prodID)
	err = utils.GetCache(cacheId, &prod)
	if err == nil { // if found in cache
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(prod)
		return
	}

	// get from db
	query := `SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`
	row := db.DB.QueryRow(query, id)
	err = row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CategoryID, &prod.ImageURL)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	_ = utils.SetCache(cacheId, prod, time.Minute*5)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(prod)
}

func UpdateOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
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
		http.Error(w, "Error updating", http.StatusBadRequest)
		return
	}

	// clear from cache
	cacheId := fmt.Sprintf("product_" + prodID)
	_ = utils.DeleteCache(cacheId)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(prod)
}

func DeleteOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM products 
				WHERE id=$1`

	_, err = db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error updating", http.StatusBadRequest)
		return
	}

	// clear from cache
	cacheId := fmt.Sprintf("product_" + prodID)
	_ = utils.DeleteCache(cacheId)

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode("Deleted Successfully")
}
