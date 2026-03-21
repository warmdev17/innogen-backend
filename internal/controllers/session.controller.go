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
// @Description Retrieve a session with its lessons (display info only)
// @Tags course
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} models.SessionResponse
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

	lessons := make([]models.LessonItem, len(session.Lessons))
	for i, l := range session.Lessons {
		lessons[i] = models.LessonItem{
			ID:         l.ID,
			Title:      l.Title,
			OrderIndex: l.OrderIndex,
		}
	}

	c.JSON(http.StatusOK, models.SessionResponse{
		ID:         session.ID,
		Title:      session.Title,
		OrderIndex: session.OrderIndex,
		Lessons:    lessons,
	})
}
