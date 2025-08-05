package models

import "time"

type Product struct {
	// Product ID
	ID int `db:"id" json:"id" example:"1" description:"Unique identifier for the product"`
	// Product Name
	Name string `db:"name" json:"name" example:"Lenovo 96" description:"Name of the product"`
	// Product Description
	Description string `db:"description" json:"description" example:"A high-performance laptop" description:"Detailed description of the product"`
	// Price of the Product
	Price float64 `db:"price" json:"price" example:"999.99" description:"Price of the product"`
	// Number of Stocks
	Stock int `db:"stock" json:"stock" example:"50" description:"Available stock for the product"`
	// ID of the category the product belongs to
	CategoryID int `db:"category_id" json:"category_id" example:"2" description:"ID of the category the product belongs to"`
	// URL of the product image
	ImageURL string `db:"image_url" json:"image_url" example:"https://example.com/image.jpg" description:"URL of the product image"`
	// Timestamp when the product was created
	CreatedAt time.Time `db:"created_at" json:"created_at" example:"2023-08-01T12:00:00Z" description:"Timestamp when the product was created"`
}

// ProductCreateRequest
// @Description A product contains details like name, description, price, stock, category, and image URL.
type CreateProductRequest struct {
	Name        string  `json:"name" example:"Nokia" description:"Name of the product"`
	Description string  `json:"description" example:"A brand-new smartphone with advanced features" description:"Detailed description of the product"`
	Price       float64 `json:"price" example:"799.99" description:"Price of the product"`
	Stock       int     `json:"stock" example:"100" description:"Available stock for the product"`
	CategoryID  int     `json:"category_id" example:"3" description:"ID of the category the product belongs to"`
	ImageURL    string  `json:"image_url" example:"https://example.com/smartphone.jpg" description:"URL of the product image"`
}

// ProductUpdateRequest
// @Description A product contains details like name, description, price, stock, category, and image URL.
type UpdateProductRequest struct {
	Name        *string  `json:"name" example:"Nokia" description:"Name of the product"`
	Description *string  `json:"description" example:"A brand-new smartphone with advanced features" description:"Detailed description of the product"`
	Price       *float64 `json:"price" example:"799.99" description:"Price of the product"`
	Stock       *int     `json:"stock" example:"100" description:"Available stock for the product"`
	CategoryID  *int     `json:"category_id" example:"3" description:"ID of the category the product belongs to"`
	ImageURL    *string  `json:"image_url" example:"https://example.com/smartphone.jpg" description:"URL of the product image"`
}

// Category Model
// @Description
type Category struct {
	ID   int    `db:"id" json:"id" example:"1" description:"Unique identifier for the category"`
	Name string `db:"name" json:"name" example:"Electronics" description:"Name of the category"`
	Slug string `db:"slug" json:"slug" example:"electronics" description:"Slugified name of the category"`
}

type CreateCategoryRequest struct {
	Name string `db:"name" json:"name" example:"Electronics" description:"Name of the category"`
	Slug string `db:"slug" json:"slug" example:"electronics" description:"Slugified name of the category"`
}
