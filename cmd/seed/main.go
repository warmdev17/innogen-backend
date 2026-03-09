package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"github.com/warmdev17/innogen-backend/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	email := "admin@admin.com"
	password := "admin123"

	// Check if user exists
	var existing models.User
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		log.Println("Admin user already exists")
		return
	}

	hash, err := services.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := models.User{
		Email:    email,
		Password: hash,
		Role:     "admin",
	}

	if err := database.DB.Create(&adminUser).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("Successfully created admin user: %s / %s\n", email, password)
}
