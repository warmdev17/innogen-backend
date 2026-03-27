package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

// GetLessons godoc
// @Summary Get all lessons (with optional filtering)
// @Description Retrieve lessons, optionally filtered by sessionId
// @Tags course
// @Accept json
// @Produce json
// @Param sessionId query int false "Filter by session ID"
// @Success 200 {array} models.Lesson
// @Router /lessons [get]
func GetLessons(c *gin.Context) {
	var lessons []models.Lesson
	query := database.DB

	// Filter by sessionId if provided
	if sessionId := c.Query("sessionId"); sessionId != "" {
		id, err := strconv.ParseUint(sessionId, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sessionId"})
			return
		}
		query = query.Where("subject_session_id = ?", id)
	}

	err := query.Order("order_index ASC").Find(&lessons).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lessons"})
		return
	}

	c.JSON(http.StatusOK, lessons)
}

// GetLesson godoc
// @Summary Get lesson by ID
// @Description Retrieve a lesson with its problems (display info only)
// @Tags course
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} models.LessonResponse
// @Failure 404 {object} map[string]string
// @Router /lessons/{id} [get]
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

	problems := make([]models.LessonProblemItem, 0, len(lesson.Problems))
	for _, lp := range lesson.Problems {
		if lp.Problem == nil {
			continue
		}
		problems = append(problems, models.LessonProblemItem{
			ID:             lp.Problem.ID,
			Slug:           lp.Problem.Slug,
			Title:          lp.Problem.Title,
			Difficulty:     lp.Problem.Difficulty,
			AcceptanceRate: lp.Problem.AcceptanceRate,
			OrderIndex:     lp.OrderIndex,
		})
	}

	c.JSON(http.StatusOK, models.LessonResponse{
		ID:         lesson.ID,
		SessionID:  lesson.SubjectSessionID,
		Title:      lesson.Title,
		ContentMd:  lesson.ContentMd,
		OrderIndex: lesson.OrderIndex,
		Problems:   problems,
	})
}
