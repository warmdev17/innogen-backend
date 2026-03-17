// Package controllers
package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
)

type SubmitRequest struct {
	ProblemID uint   `json:"problemId"`
	Code      string `json:"code"`
	Language  string `json:"language"`
}

func Submit(c *gin.Context) {
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userIDFloat := c.GetFloat64("user_id")
	userID := uint(userIDFloat)

	sub := models.Submission{
		ID:        uuid.New(),
		UserID:    userID,
		ProblemID: req.ProblemID,
		Code:      req.Code,
		Language:  req.Language,
		Status:    "pending",
	}

	database.DB.Create(&sub)

	jobData := map[string]any{
		"submission_id": sub.ID,
		"problem_id":    sub.ProblemID,
	}

	jsonData, _ := json.Marshal(jobData)

	err := database.Rdb.RPush(database.Ctx, "judge_queue", jsonData).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to queue submission"})
		return
	}

	c.JSON(201, gin.H{
		"message":    "Submission queued",
		"submission": sub,
	})
}

func GetSubmission(c *gin.Context) {
	id := c.Param("id")
	var sub models.Submission

	if err := database.DB.Where("id = ?", id).First(&sub).Error; err != nil {
		c.JSON(404, gin.H{"error": "Submission not found"})
		return
	}

	// Prevent users from seeing other people's code unless they are admin/teacher. 
	// For now, simpler implementation: just check if user_id matches
	userIDFloat := c.GetFloat64("user_id")
	userID := uint(userIDFloat)
	role := c.GetString("role")
	
	if sub.UserID != userID && role != "admin" && role != "teacher" {
		c.JSON(403, gin.H{"error": "You do not have permission to view this submission"})
		return
	}

	c.JSON(200, sub)
}
