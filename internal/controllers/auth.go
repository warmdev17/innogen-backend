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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"github.com/warmdev17/innogen-backend/internal/services"
	"github.com/warmdev17/innogen-backend/internal/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type MeResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FullName  string `json:"fullName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	if err := services.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate access token (15 minutes)
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate access token"})
		return
	}

	// Generate refresh token (30 days) and store in database
	refreshTokenStr, err := services.CreateRefreshToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshTokenStr,
	})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get information about the authenticated user
// @Tags me
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MeResponse
// @Failure 401 {object} ErrorResponse
// @Router /me [get]
func GetCurrentUser(c *gin.Context) {
	userIdFloat := c.GetFloat64("user_id")
	userID := uint(userIdFloat)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"username":  user.Username,
		"fullName":  user.FullName,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param RefreshTokenRequest body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(401, gin.H{"error": "refresh token is required"})
		return
	}

	// Parse the refresh token
	refreshToken, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	// Extract claims
	claims, err := utils.GetTokenClaims(refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid token claims"})
		return
	}

	// Verify token is not expired
	if float64(claims["exp"].(float64)) < float64(time.Now().Unix()) {
		c.JSON(401, gin.H{"error": "refresh token expired"})
		return
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	// Verify refresh token exists in database and is not revoked
	_, err = services.VerifyRefreshToken(req.RefreshToken, userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "refresh token not found or revoked"})
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(userID, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate access token"})
		return
	}

	// Rotate refresh token for security (revoke old, create new)
	newRefreshToken, err := services.RotateRefreshToken(req.RefreshToken, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to rotate refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Revoke the refresh token and logout the user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	// Get refresh token from request body
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(400, gin.H{"error": "refresh token is required"})
		return
	}

	// Get user ID from access token
	userIdFloat := c.GetFloat64("user_id")
	userID := uint(userIdFloat)

	// Revoke the refresh token
	err := services.RevokeRefreshToken(req.RefreshToken, userID)
	if err != nil {
		// Token might not exist or already revoked, but still return success
		c.JSON(200, gin.H{"message": "logged out successfully"})
		return
	}

	c.JSON(200, gin.H{"message": "logged out successfully"})
}

// LogoutAll godoc
// @Summary Logout from all devices
// @Description Revoke all refresh tokens for the user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout-all [post]
func LogoutAll(c *gin.Context) {
	// Get user ID from access token
	userIdFloat := c.GetFloat64("user_id")
	userID := uint(userIdFloat)

	// Revoke all refresh tokens for the user
	err := services.RevokeAllUserTokens(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(200, gin.H{"message": "logged out from all devices successfully"})
}
