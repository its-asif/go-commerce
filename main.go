package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/config"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/routes"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Welcome to GO-Commerce")
	config.LoadEnv()
	db.ConnectDB()
	db.ConnectRedis()

	router := mux.NewRouter()
	routes.GetRoutes(router)

	port := config.GetEnv("PORT")
	fmt.Println("running on localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
