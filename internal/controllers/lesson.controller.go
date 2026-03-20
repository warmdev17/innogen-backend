package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetLesson retrieves a specific lesson by ID, including its associated problems
func GetLesson(c *gin.Context) {
	id := c.Param("id")
	var lesson models.Lesson

	err := database.DB.
		Preload("Problems", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Problems.Problem").
		Where("id = ?", id).
		First(&lesson).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}
