package db

import (
	"fmt"
	"github.com/its-asif/go-commerce/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var DB *sqlx.DB

func ConnectDB() {
	var err error
	DB, err = sqlx.Connect("postgres", config.GetEnv("DB_URL"))
	if err != nil {
		log.Fatal("Error connecting DB")
	}
	fmt.Println("DB connected successfully")

	//defer DB.Close()
}
