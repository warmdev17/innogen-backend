// Package controllers
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
)

type CreateProblemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	TimeLimitMs int    `json:"time_limit_ms"`
	MemoryLimit int    `json:"memory_limit_mb"`
}

func CreateProblem(c *gin.Context) {
	var req CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	problem := models.Problem{
		Title:         req.Title,
		Description:   req.Description,
		Difficulty:    req.Difficulty,
		TimeLimitMs:   req.TimeLimitMs,
		MemoryLimitMB: req.MemoryLimit,
		CreatedBy:     userID,
	}

	database.DB.Create(&problem)

	c.JSON(http.StatusCreated, problem)
}

func GetProblems(c *gin.Context) {
	var problems []models.Problem
	database.DB.Find(&problems)

	c.JSON(200, problems)
}

func GetProblemByID(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem

	if err := database.DB.First(&problem, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Problem not found"})
		return
	}

	c.JSON(200, problem)
}
