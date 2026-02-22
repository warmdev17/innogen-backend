// Package controllers
package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
)

type SubmitRequest struct {
	ProblemID uint   `json:"problem_id"`
	Code      string `json:"code"`
	Language  string `json:"language"`
}

func Submit(c *gin.Context) {
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	sub := models.Submission{
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
