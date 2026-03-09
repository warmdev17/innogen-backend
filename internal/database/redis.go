// Package database
package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb redis.Client

var Ctx = context.Background()

func ConnectRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	Rdb = *redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
}
