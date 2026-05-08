package handlers

import (
	"encoding/json"
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

	// db query
	err = db.CreateProduct(input)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	// invalid redis cache
	_ = utils.DeleteCache("all_products")

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode("Product created successfully")
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
	prod, err = db.GetAllProduct()

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
	cacheId := "product_" + prodID
	err = utils.GetCache(cacheId, &prod)
	if err == nil { // if found in cache
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(prod)
		return
	}

	// get from db
	prod, err = db.GetSingleProduct(id)

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

	// update from db
	err = db.UpdateProduct(id, input, w)

	if err != nil {
		http.Error(w, "Error updating", http.StatusBadRequest)
		return
	}

	// clear from cache
	cacheId := "product_" + prodID
	_ = utils.DeleteCache(cacheId)
	// invalidate all_products cache
	_ = utils.DeleteCache("all_products")

	w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(input)
	w.Write([]byte("Product updated successfully"))
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

	// delete from db
	err = db.DeleteProduct(id, w)

	if err != nil {
		http.Error(w, "Error updating", http.StatusBadRequest)
		return
	}

	// clear from cache
	cacheId := "product_" + prodID
	_ = utils.DeleteCache(cacheId)

	// invalidate all_products cache
	_ = utils.DeleteCache("all_products")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Deleted Successfully")
}
