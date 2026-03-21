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
// @Description Retrieve list of all subjects for display
// @Tags course
// @Accept json
// @Produce json
// @Success 200 {array} models.SubjectListItem
// @Router /subjects [get]
func GetSubjects(c *gin.Context) {
	var subjects []models.Subject
	if err := database.DB.Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
		return
	}

	result := make([]models.SubjectListItem, len(subjects))
	for i, s := range subjects {
		result[i] = models.SubjectListItem{
			ID:          s.ID,
			Name:        s.Name,
			Slug:        s.Slug,
			Description: s.Description,
			Color:       s.Color,
		}
	}
	c.JSON(http.StatusOK, result)
}

// GetSubject godoc
// @Summary Get subject by slug
// @Description Retrieve a subject with its sessions and lessons (display info only)
// @Tags course
// @Accept json
// @Produce json
// @Param slug path string true "Subject Slug"
// @Success 200 {object} models.SubjectDetailResponse
// @Failure 404 {object} map[string]string
// @Router /subjects/{slug} [get]
func GetSubject(c *gin.Context) {
	slug := c.Param("slug")
	var subject models.Subject

	err := database.DB.
		Preload("Sessions", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
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

	sessions := make([]models.SessionItem, len(subject.Sessions))
	for i, sess := range subject.Sessions {
		lessons := make([]models.LessonItem, len(sess.Lessons))
		for j, l := range sess.Lessons {
			lessons[j] = models.LessonItem{
				ID:         l.ID,
				Title:      l.Title,
				OrderIndex: l.OrderIndex,
			}
		}
		sessions[i] = models.SessionItem{
			ID:         sess.ID,
			Title:      sess.Title,
			OrderIndex: sess.OrderIndex,
			Lessons:    lessons,
		}
	}

	c.JSON(http.StatusOK, models.SubjectDetailResponse{
		ID:          subject.ID,
		Name:        subject.Name,
		Slug:        subject.Slug,
		Description: subject.Description,
		Color:       subject.Color,
		Sessions:    sessions,
	})
}
