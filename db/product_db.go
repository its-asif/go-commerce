package db

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/its-asif/go-commerce/models"
)

var GetAllProduct = func() ([]models.Product, error) {
	var prod []models.Product
	query := `SELECT * FROM products`
	err := DB.Select(&prod, query)
	return prod, err
}

var GetSingleProduct = func(id int) (models.Product, error) {
	var prod models.Product
	query := `SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`
	row := DB.QueryRow(query, id)
	err := row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CategoryID, &prod.ImageURL)
	return prod, err
}

var CreateProduct = func(input models.Product) error {
	query := `INSERT INTO products (name, description, price, stock, category_id, image_url)
				VALUES($1,$2,$3,$4,$5,$6)
				RETURNING id, created_at`
	err := DB.QueryRowx(query, input.Name, input.Description, input.Price, input.Stock, input.CategoryID, input.ImageURL).Scan(&input.ID, &input.CreatedAt)
	return err
}

var UpdateProduct = func(id int, input models.UpdateProductRequest, w http.ResponseWriter) error {
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
		return nil
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
	_, err := DB.Exec(query, args...)
	return err
}

var joinStrings = func(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

var DeleteProduct = func(id int, w http.ResponseWriter) error {
	var prod models.Product
	query := `SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`
	productRow := DB.QueryRow(`SELECT id, name, description, price, stock, category_id, image_url FROM products WHERE id=$1`, id)
	err := productRow.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CategoryID, &prod.ImageURL)
	if err != nil {
		_ = json.NewEncoder(w).Encode("the product doesn't exist")
		return fmt.Errorf("the product doesn't exist")
	}
	query = `DELETE FROM products 
				WHERE id=$1`
	_, err = DB.Exec(query, id)
	return err
}
