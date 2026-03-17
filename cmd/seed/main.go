package main

import (
	"log"
	"os"

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
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is not set")
	}

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
	
	// Create sample problem
	problem := models.Problem{
		Title:          "Two Sum",
		Slug:           "two-sum",
		ProblemMd:      "Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to `target`.",
		Difficulty:     "Easy",
		TimeLimitMs:    1000,
		MemoryLimitKb:  256000,
		IsPublished:    true,
		CreatedBy:      adminUser.ID,
		AcceptanceRate: 0.0,
	}
	
	if err := database.DB.Create(&problem).Error; err != nil {
		log.Fatalf("Failed to create basic problem: %v", err)
	}
	log.Println("Successfully created sample problem: Two Sum")
	
	// Create test cases
	testCases := []models.Testcase{
		{
			ProblemID:      problem.ID,
			InputData:      "2\n7\n11\n15\n9",
			ExpectedOutput: "0\n1",
			IsHidden:       false,
			Role:           "sample",
		},
		{
			ProblemID:      problem.ID,
			InputData:      "3\n2\n4\n6",
			ExpectedOutput: "1\n2",
			IsHidden:       true,
			Role:           "hidden",
		},
		{
			ProblemID:      problem.ID,
			InputData:      "3\n3\n6",
			ExpectedOutput: "0\n1",
			IsHidden:       true,
			Role:           "edge_case",
		},
	}
	
	for _, tc := range testCases {
		if err := database.DB.Create(&tc).Error; err != nil {
			log.Fatalf("Failed to create test case: %v", err)
		}
	}
	
	log.Println("Successfully created 3 sample test cases")
}
