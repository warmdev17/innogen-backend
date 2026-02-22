// Package controllers
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
)

func CreateTestcase(c *gin.Context) {
	var tc models.Testcase
	if err := c.ShouldBindJSON(&tc); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&tc)
	c.JSON(201, tc)
}
