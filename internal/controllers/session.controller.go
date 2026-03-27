package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetSessions godoc
// @Summary Get all sessions (with optional filtering)
// @Description Retrieve sessions, optionally filtered by subjectId
// @Tags course
// @Accept json
// @Produce json
// @Param subjectId query int false "Filter by subject ID"
// @Success 200 {array} models.SubjectSession
// @Router /sessions [get]
func GetSessions(c *gin.Context) {
	var sessions []models.SubjectSession
	query := database.DB

	// Filter by subjectId if provided
	if subjectId := c.Query("subjectId"); subjectId != "" {
		id, err := strconv.ParseUint(subjectId, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subjectId"})
			return
		}
		query = query.Where("subject_id = ?", id)
	}

	err := query.Order("order_index ASC").Find(&sessions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

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
			SessionID:  l.SubjectSessionID,
			Title:      l.Title,
			ContentMd:  l.ContentMd,
			OrderIndex: l.OrderIndex,
		}
	}

	c.JSON(http.StatusOK, models.SessionResponse{
		ID:         session.ID,
		SubjectID:  session.SubjectID,
		Title:      session.Title,
		OrderIndex: session.OrderIndex,
		Lessons:    lessons,
	})
}
