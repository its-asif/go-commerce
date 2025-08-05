package db

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/its-asif/go-commerce/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

// Redis connection
var Rdb *redis.Client

func ConnectRedis() {
	dbNum, err := strconv.Atoi(config.GetEnv("REDIS_DB"))

	if err != nil {
		dbNum = 0
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_URL")[8:],
		Password: config.GetEnv("REDIS_PASSWORD"),
		DB:       dbNum,
	})

	//	test connection
	ctx := context.Background()
	_, err = Rdb.Ping(ctx).Result()

	if err != nil {
		log.Println("Error connecting redis", err)
	}
	fmt.Println("Redis connected")
}
