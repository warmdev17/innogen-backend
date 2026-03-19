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

	// Create admin user
	adminEmail := "admin@admin.com"
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is not set")
	}

	createUser(adminEmail, "admin", "Admin User", adminPassword, "admin")

	// Create additional test users
	testUsers := []struct {
		email    string
		username string
		fullName string
		role     string
		password string
	}{
		{"teacher@innogen.com", "teacher", "John Teacher", "teacher", "teacher123"},
		{"student1@innogen.com", "student1", "Alice Student", "student", "student123"},
		{"student2@innogen.com", "student2", "Bob Student", "student", "student123"},
		{"student3@innogen.com", "student3", "Charlie Student", "student", "student123"},
	}

	for _, u := range testUsers {
		createUser(u.email, u.username, u.fullName, u.password, u.role)
	}

	// Create sample problem using the admin user
	var adminUser models.User
	if err := database.DB.Where("email = ?", adminEmail).First(&adminUser).Error; err != nil {
		log.Fatalf("Failed to find admin user: %v", err)
	}

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
	log.Println("\n=== Test Accounts Created ===")
	log.Println("Admin: admin@admin.com / " + adminPassword)
	log.Println("Teacher: teacher@innogen.com / teacher123")
	log.Println("Student 1: student1@innogen.com / student123")
	log.Println("Student 2: student2@innogen.com / student123")
	log.Println("Student 3: student3@innogen.com / student123")
	log.Println("===============================\n")
}

func createUser(email, username, fullName, password, role string) {
	var existing models.User
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		log.Printf("User %s already exists, skipping\n", email)
		return
	}

	hash, err := services.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password for %s: %v", email, err)
	}

	user := models.User{
		Email:    email,
		Username: username,
		FullName: fullName,
		Password: hash,
		Role:     role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user %s: %v", email, err)
	}

	log.Printf("Created user: %s (role: %s)\n", email, role)
}
