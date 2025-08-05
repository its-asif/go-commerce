package routes

import (
	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/handlers"
	"github.com/its-asif/go-commerce/middleware"
)

func GetRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	//	user auth
	r.HandleFunc("/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/auth/login", handlers.Login).Methods("POST")

	//	auth middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	adminRoute := api.NewRoute().Subrouter()
	adminRoute.Use(middleware.AdminMiddleware)

	//	Products
	adminRoute.HandleFunc("/products", handlers.CreateProducts).Methods("POST")
	api.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	api.HandleFunc("/products/{id}", handlers.GetOneProduct).Methods("GET")
	api.HandleFunc("/products/{id}", handlers.UpdateOneProduct).Methods("PUT")
	adminRoute.HandleFunc("/products/{id}", handlers.DeleteOneProduct).Methods("DELETE")

	// cart
	api.HandleFunc("/cart", handlers.GetCarts).Methods("GET")
	api.HandleFunc("/cart", handlers.AddToCart).Methods("POST")
	api.HandleFunc("/cart/{product_id}", handlers.RemoveFromCart).Methods("DELETE")

	// orders
	api.HandleFunc("/orders/checkout", handlers.Checkout).Methods("POST")
	api.HandleFunc("/orders", handlers.GetOrders).Methods("GET")

	// category
	adminRoute.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")
	adminRoute.HandleFunc("/categories", handlers.GetAllCategories).Methods("GET")
}
