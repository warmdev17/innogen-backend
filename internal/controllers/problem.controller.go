// Package controllers
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"gorm.io/gorm"
)

type TestcaseRequest struct {
	InputData      string `json:"inputData" binding:"required"`
	ExpectedOutput string `json:"expectedOutput" binding:"required"`
	IsHidden       bool   `json:"isHidden"`
}

type CreateProblemRequest struct {
	Title         string            `json:"title" binding:"required"`
	Slug          string            `json:"slug" binding:"required"`
	ProblemMd     string            `json:"problemMd"`
	Difficulty    string            `json:"difficulty"`
	TimeLimitMs   int               `json:"timeLimitMs"`
	MemoryLimitKb int               `json:"memoryLimitKb"`
	Testcases     []TestcaseRequest `json:"testcases" binding:"required"`
}

type ProblemListResponse struct {
	ID         uint   `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
}

type ProblemDetailResponse struct {
	ID             uint      `json:"id"`
	Slug           string    `json:"slug"`
	AcceptanceRate float64   `json:"acceptanceRate"`
	Title          string    `json:"title"`
	Difficulty     string    `json:"difficulty"`
	ProblemMd      string    `json:"problemMd"`
	TimeLimitMs    int       `json:"timeLimitMs"`
	MemoryLimitKb  int       `json:"memoryLimitKb"`
	IsPublished    bool      `json:"isPublished"`
}

func CreateProblem(c *gin.Context) {
	var req CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdFloat := c.GetFloat64("user_id")
	userID := uint(userIdFloat)

	var testcases []models.Testcase
	for _, tc := range req.Testcases {
		testcases = append(testcases, models.Testcase{
			InputData:      tc.InputData,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
		})
	}

	problem := models.Problem{
		Title:         req.Title,
		Slug:          req.Slug,
		ProblemMd:     req.ProblemMd,
		Difficulty:    req.Difficulty,
		TimeLimitMs:   req.TimeLimitMs,
		MemoryLimitKb: req.MemoryLimitKb,
		CreatedBy:     userID,
		Testcases:     testcases,
	}

	if err := database.DB.Create(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create problem"})
		return
	}

	c.JSON(http.StatusCreated, problem)
}

func GetProblems(c *gin.Context) {
	var problems []models.Problem
	database.DB.Select("id, slug, title, difficulty").Find(&problems)

	var response []ProblemListResponse
	for _, p := range problems {
		response = append(response, ProblemListResponse{
			ID:         p.ID,
			Slug:       p.Slug,
			Title:      p.Title,
			Difficulty: p.Difficulty,
		})
	}

	// Make sure we return an empty array [] instead of null if no problems
	if response == nil {
		response = []ProblemListResponse{}
	}

	c.JSON(200, response)
}

func GetProblemByID(c *gin.Context) {
	id := c.Param("id")
	var p models.Problem

	if err := database.DB.First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Problem not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	response := ProblemDetailResponse{
		ID:             p.ID,
		Slug:           p.Slug,
		AcceptanceRate: p.AcceptanceRate,
		Title:          p.Title,
		Difficulty:     p.Difficulty,
		ProblemMd:      p.ProblemMd,
		TimeLimitMs:    p.TimeLimitMs,
		MemoryLimitKb:  p.MemoryLimitKb,
		IsPublished:    p.IsPublished,
	}

	c.JSON(200, response)
}
