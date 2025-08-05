package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
)

// @Summary		Create Product
// @Description	Create a new product
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string			true	"Bearer + JWT_Token"
// @Param			input			body		models.CreateProductRequest	true	"Product Input"
// @Success		201				{object}	models.Product
// @Failure		400				{string}	string	"Invalid Input"
// @Failure		402				{string}	string	"Unauthorized"
// @Failure		500				{string}	string	"Server Error"
// @Router			/api/products [POST]
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

// @Summary		Get All Products
// @Description	Get all the products from db
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Bearer + JWT_Token"
// @Success		200				{object}	models.Product
// @Failure		400				{string}	string	"Bad Request"
// @Router			/api/products [GET]
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

// @Summary		Get One Products
// @Description	Get a single product by ID
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Bearer + JWT_Token"
// @Param			id	path		int	true	"Product ID"
// @Success		200				{object}	models.Product
// @Failure		400				{string}	string	"Bad Request"
// @Router			/api/products/{id} [GET]
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

// @Summary		Update Product
// @Description	Update a product by ID
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string			true	"Bearer + JWT_Token"
// @Param			id	path		int	true	"Product ID"
// @Param			input			body		models.UpdateProductRequest	true	"Product Input"
// @Success		200				{object}	models.Product
// @Failure		400				{string}	string	"Bad request"
// @Router			/api/products/{id} [PUT]
func UpdateOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	//var input models.Product
	var input models.UpdateProductRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}

	updates := []string{}
	args := []interface{}{}
	argPosition := 1

	if input.Name != nil {
		updates = append(updates, fmt.Sprintf("name=$%d", argPosition))
		args = append(args, *input.Name)
		argPosition++
	}
	if input.Description != nil {
		updates = append(updates, fmt.Sprintf("description=$%d", argPosition))
		args = append(args, *input.Description)
		argPosition++
	}
	if input.Price != nil {
		updates = append(updates, fmt.Sprintf("price=$%d", argPosition))
		args = append(args, *input.Price)
		argPosition++
	}
	if input.Stock != nil {
		updates = append(updates, fmt.Sprintf("stock=$%d", argPosition))
		args = append(args, *input.Stock)
		argPosition++
	}
	if input.CategoryID != nil {
		updates = append(updates, fmt.Sprintf("category_id=$%d", argPosition))
		args = append(args, *input.CategoryID)
		argPosition++
	}
	if input.ImageURL != nil {
		updates = append(updates, fmt.Sprintf("image_url=$%d", argPosition))
		args = append(args, *input.ImageURL)
		argPosition++
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	//var result models.Product
	//query := `UPDATE products
	//			SET name=$1, description=$2, price=$3, stock=$4, category_id=$5, image_url=$6
	//			WHERE id=$7`
	//

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id=$%d`,
		joinStrings(updates, ", "), argPosition)
	args = append(args, id)

	//_, err = db.DB.Exec(query, prod.Name, prod.Description, prod.Price, prod.Stock, prod.CategoryID, prod.ImageURL, id)
	_, err = db.DB.Exec(query, args...)

	if err != nil {
		http.Error(w, "Error updating", http.StatusBadRequest)
		return
	}

	// clear from cache
	cacheId := fmt.Sprintf("product_" + prodID)
	_ = utils.DeleteCache(cacheId)
	// invalidate all_products cache
	_ = utils.DeleteCache("all_products")

	w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(input)
	w.Write([]byte("Product updated successfully"))
}

func joinStrings(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// @Summary		Delete Product
// @Description	Delete a product by ID
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string			true	"Bearer + JWT_Token"
// @Param			id	path		int	true	"Product ID"
// @Success		200				{object}	models.Product
// @Failure		400				{string}	string	"Bad request"
// @Router			/api/products/{id} [DELETE]
func DeleteOneProduct(w http.ResponseWriter, r *http.Request) {
	prodID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodID)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	var prod models.Product
	//productRow := db.DB.QueryRow("SELECT * WHERE id $1", id)
	query := `SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`
	productRow := db.DB.QueryRow(`SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`, id)
	err = productRow.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CategoryID, &prod.ImageURL)

	// check if it exists
	if err != nil {
		_ = json.NewEncoder(w).Encode("the product doesn't exist")
		return
	}

	query = `DELETE FROM products 
				WHERE id=$1`

	_, err = db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error updating", http.StatusBadRequest)
		return
	}

	// clear from cache
	cacheId := fmt.Sprintf("product_" + prodID)
	_ = utils.DeleteCache(cacheId)

	// invalidate all_products cache
	_ = utils.DeleteCache("all_products")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Deleted Successfully")
}
