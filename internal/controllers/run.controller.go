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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/services"
)

type RunRequest struct {
	ProblemID uint   `json:"problemId"` // optional if just running bare code, but good for saving
	Code      string `json:"code" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Stdin     string `json:"stdin"`
}

// RunCode godoc
// @Summary Run code directly
// @Description Execute code without test cases and return the result
// @Tags run
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RunRequest true "Code execution details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /run [post]
func RunCode(c *gin.Context) {
	var req RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdFloat, exists := c.Get("user_id")
	if exists {
		userID := uint(userIdFloat.(float64))
        if req.ProblemID > 0 {
            redisKey := fmt.Sprintf("run_code:%d:%d", userID, req.ProblemID)
            database.Rdb.Set(database.Ctx, redisKey, req.Code, 0)
        }
	}

	resp, err := services.RunCodeWithInput(req.Code, req.Language, req.Stdin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute code against Piston"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
