// Package database
package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Rdb redis.Client

var Ctx = context.Background()

func ConnectRedis() {
	Rdb = *redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
}
