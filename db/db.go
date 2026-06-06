package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
		log.Fatal("Error connecting DB: ", err)
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

	// Prefer REDIS_URL if provided, e.g., redis://:password@host:6379/0
	redisURL := strings.TrimSpace(config.GetEnv("REDIS_URL"))
	if redisURL != "" {
		opt, perr := redis.ParseURL(redisURL)
		if perr == nil {
			// Override DB if REDIS_DB is set explicitly
			opt.DB = dbNum
			if pw := strings.TrimSpace(config.GetEnv("REDIS_PASSWORD")); pw != "" {
				opt.Password = pw
			}
			Rdb = redis.NewClient(opt)
		}
	}

	// Fallback to options if URL parse failed or URL not provided
	if Rdb == nil {
		Rdb = redis.NewClient(&redis.Options{
			Addr:     strings.TrimPrefix(redisURL, "redis://"),
			Password: config.GetEnv("REDIS_PASSWORD"),
			DB:       dbNum,
		})
	}

	//	test connection
	ctx := context.Background()
	_, err = Rdb.Ping(ctx).Result()

	if err != nil {
		log.Println("Error connecting redis", err)
	}
	fmt.Println("Redis connected")
}
