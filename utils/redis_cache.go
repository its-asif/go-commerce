package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/its-asif/go-commerce/db"
)

var GetCache = func(key string, destination interface{}) error {
	ctx := context.Background()
	val, err := db.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), destination)
}

var SetCache = func(key string, value interface{}, expiration time.Duration) error {
	fmt.Println("Setting cache for key:", key)
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = db.Rdb.Set(ctx, key, jsonValue, expiration).Err()
	return err
}

var DeleteCache = func(key string) error {
	ctx := context.Background()
	err := db.Rdb.Del(ctx, key).Err()
	return err
}

var GetCacheString = func(key string, destination interface{}) (string, error) {
	ctx := context.Background()
	return db.Rdb.Get(ctx, key).Result()
}

var SetCacheString = func(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return db.Rdb.Set(ctx, key, value, expiration).Err()
}
