package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetSubjects godoc
// @Summary Get all subjects
// @Description Retrieve list of all subjects (course structures)
// @Tags course
// @Accept json
// @Produce json
// @Success 200 {array} models.Subject
// @Router /subjects [get]
func GetSubjects(c *gin.Context) {
	var subjects []models.Subject
	// Basic info only, don't load all nested data for the list view
	if err := database.DB.Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
		return
	}
	c.JSON(http.StatusOK, subjects)
}

// GetSubject godoc
// @Summary Get subject by slug
// @Description Retrieve a specific subject by its slug, including nested sessions and lessons
// @Tags course
// @Accept json
// @Produce json
// @Param slug path string true "Subject Slug"
// @Success 200 {object} models.Subject
// @Failure 404 {object} map[string]string
// @Router /subjects/{slug} [get]
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
