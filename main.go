package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/config"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/routes"

	_ "github.com/its-asif/go-commerce/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			GO Commerce API
//	@version		1.0
//	@description	This is a RESTful API for an E-commerce platform built with Go.
//	@host			localhost:8000
//	@BasePath		/

func main() {
	fmt.Println("Welcome to GO-Commerce")
	config.LoadEnv()
	db.ConnectDB()
	db.ConnectRedis()

	router := mux.NewRouter()
	routes.GetRoutes(router)

	// Swagger UI route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := config.GetEnv("PORT")
	fmt.Println("running on localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
