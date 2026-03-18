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
	"log"

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
	Token string `json:"token"`
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

// SendOTP godoc
// @Summary Send OTP for registration
// @Description Sends a 6-digit OTP to the user's email for registration (Currently disabled)
// @Tags auth
// @Accept json
// @Produce json
// @Param email body object{email=string} true "User email"
// @Success 200 {object} object{message=string} "OTP sent successfully"
// @Failure 400 {object} object{error=string} "Bad request"
// @Router /auth/send-otp [post]
func SendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "valid email is required"})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "user exists"})
		return
	}

	otp := services.GenerateOTP()
	if err := services.StoreOTP(req.Email, otp); err != nil {
		c.JSON(500, gin.H{"error": "failed to store OTP"})
		return
	}

	go func() {
		err := services.SendEmail([]string{req.Email}, "Innogen Registration OTP", "Your OTP is: "+otp+"\n\nIt expires in 5 minutes.")
		if err != nil {
			log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
		} else {
			log.Printf("OTP email successfully sent to %s", req.Email)
		}
	}()

	c.JSON(200, gin.H{"message": "OTP sent"})
}

// Register godoc
// @Summary Register a new user
// @Description Verifies the OTP and creates a new user with the student role (Currently disabled)
// @Tags auth
// @Accept json
// @Produce json
// @Param email body object{email=string} true "User email"
// @Param password body object{password=string} true "User password"
// @Param otp body object{otp=string} true "6-digit OTP"
// @Success 201 {object} object{message=string} "User registered successfully"
// @Failure 400 {object} object{error=string} "Bad request"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OTP      string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !services.VerifyOTP(req.Email, req.OTP) {
		c.JSON(400, gin.H{"error": "invalid or expired OTP"})
		return
	}

	hash, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "hash failed"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: hash,
		Role:     "student",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "user exists"})
		return
	}

	c.JSON(201, gin.H{"message": "registered"})
}

// Login godoc
// @Summary User login
// @Description Authenticates the user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param email body object{email=string} true "User email"
// @Param password body object{password=string} true "User password"
// @Success 200 {object} object{token=string} "Login successful"
// @Failure 401 {object} object{error=string} "Invalid credentials"
// @Router /auth/login [post]
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
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	token, _ := utils.GenerateToken(user.ID, user.Role)

	c.JSON(200, gin.H{"token": token})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Returns the current authenticated user's information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{id=uint,email=string,username=string,fullName=string,createdAt=string,updatedAt=string} "User details"
// @Failure 404 {object} object{error=string} "User not found"
// @Router /me [get]
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
