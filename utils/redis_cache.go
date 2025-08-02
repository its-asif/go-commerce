package utils

import (
	"context"
	"encoding/json"
	"github.com/its-asif/go-commerce/db"
	"time"
)

func GetCache(key string, destination interface{}) error {
	ctx := context.Background()
	val, err := db.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), destination)
}

func SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = db.Rdb.Set(ctx, key, jsonValue, expiration).Err()
	return err
}

func DeleteCache(key string) error {
	ctx := context.Background()
	err := db.Rdb.Del(ctx, key).Err()
	return err
}

func GetCacheString(key string, destination interface{}) (string, error) {
	ctx := context.Background()
	return db.Rdb.Get(ctx, key).Result()
}

func SetCacheString(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return db.Rdb.Set(ctx, key, value, expiration).Err()
}
