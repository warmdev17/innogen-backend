package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/judge"
	"github.com/warmdev17/innogen-backend/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using system environment variables)")
	}

	database.Connect()
	database.ConnectRedis()

	go judge.StartWorker()

	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Println("Backend running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
