// @title Innogen Backend API
// @version 1.0
// @description API for competitive programming platform
// @host code.innogenlab.com
// @BasePath /api
// @schemes https http
package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/judge"
	"github.com/warmdev17/innogen-backend/internal/routes"

	_ "github.com/warmdev17/innogen-backend/docs"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using system environment variables)")
	}

	database.Connect()
	database.ConnectRedis()

	go judge.StartWorker(5)

	r := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	log.Println("Backend running on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
