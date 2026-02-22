// Package database
package database

import (
	"log"

	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=innogen password=maiphuongdangyeu dbname=innogendb port=5433 sslmode=disable"

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
