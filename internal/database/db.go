// Package database
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	portStr := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || portStr == "" {
		log.Fatal("One or more database environment variables are missing (POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT)")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, portStr)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	DB = db
	// AutoMigrate all models
	if err := DB.AutoMigrate(
		&models.Subject{},
		&models.SubjectSession{},
		&models.Lesson{},
		&models.LessonProblem{},
		&models.User{},
		&models.Problem{},
		&models.Testcase{},
		&models.Submission{},
		&models.Tag{},
		&models.RefreshToken{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database connected")
}
