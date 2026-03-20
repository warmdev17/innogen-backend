package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetSubjects retrieves all subjects
func GetSubjects(c *gin.Context) {
	var subjects []models.Subject
	// Basic info only, don't load all nested data for the list view
	if err := database.DB.Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
		return
	}
	c.JSON(http.StatusOK, subjects)
}

// GetSubject retrieves a specific subject by slug with complete hierarchy
func GetSubject(c *gin.Context) {
	slug := c.Param("slug")
	var subject models.Subject

	err := database.DB.
		Preload("Sessions", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC") // Ensure ascending order
		}).
		Preload("Sessions.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Where("slug = ?", slug).
		First(&subject).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	c.JSON(http.StatusOK, subject)
}
