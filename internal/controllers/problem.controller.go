// @title Innogen Backend API
// @version 1.0
// @description API for competitive programming platform
// @host code.innogenlab.com
// @BasePath /api
// @schemes https http
// @securityDefinitions.apiKey BearerAuth
// @type apiKey
// @in header
// @name Authorization
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

// CreateProblem godoc
// @Summary Create a new problem
// @Description Create a new problem with test cases (admin/teacher only)
// @Tags problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProblemRequest true "Problem details"
// @Success 201 {object} ProblemDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/problems [post]
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

// GetProblems godoc
// @Summary Get all problems
// @Description Retrieve list of all problems
// @Tags problems
// @Accept json
// @Produce json
// @Success 200 {array} ProblemListResponse
// @Router /problems [get]
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

// GetProblemByID godoc
// @Summary Get problem by ID
// @Description Retrieve a specific problem by its ID
// @Tags problems
// @Accept json
// @Produce json
// @Param id path int true "Problem ID"
// @Success 200 {object} ProblemDetailResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /problems/{id} [get]
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
