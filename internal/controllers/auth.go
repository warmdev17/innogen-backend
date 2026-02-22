// Package controllers
package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"github.com/warmdev17/innogen-backend/internal/services"
	"github.com/warmdev17/innogen-backend/internal/utils"
)

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
