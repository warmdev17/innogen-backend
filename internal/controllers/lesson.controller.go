package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

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
		Title:      lesson.Title,
		OrderIndex: lesson.OrderIndex,
		Problems:   problems,
	})
}
