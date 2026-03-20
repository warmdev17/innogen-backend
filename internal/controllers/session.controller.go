package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetSession retrieves a specific session by ID
func GetSession(c *gin.Context) {
	id := c.Param("id")
	var session models.SubjectSession

	err := database.DB.
		Preload("Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Where("id = ?", id).
		First(&session).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}
