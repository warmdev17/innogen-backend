// Package database
package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb redis.Client

var Ctx = context.Background()

func ConnectRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("REDIS_HOST is not set")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("REDIS_PORT is not set")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	Rdb = *redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
}
