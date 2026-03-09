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
	if host == "" {
		host = "localhost"
	}

	dsn := fmt.Sprintf("host=%s user=innogen password=maiphuongdangyeu dbname=innogendb port=5432 sslmode=disable", host)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	DB = db
	if err := DB.AutoMigrate(&models.User{}, &models.Problem{}, &models.Testcase{}, &models.Submission{}); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database connected")
}
