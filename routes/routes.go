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
}
