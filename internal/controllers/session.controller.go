package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetSession godoc
// @Summary Get session by ID
// @Description Retrieve a specific session by its ID, including nested lessons
// @Tags course
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} models.SubjectSession
// @Failure 404 {object} map[string]string
// @Router /sessions/{id} [get]
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
