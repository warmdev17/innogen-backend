// Package routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/controllers"
	"github.com/warmdev17/innogen-backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/auth/send-otp", controllers.SendOTP)
		api.POST("/auth/register", controllers.Register)
		api.POST("/auth/login", controllers.Login)

		protected := api.Group("/me")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"user_id": c.GetFloat64("user_id"),
					"role":    c.GetString("role"),
				})
			})
		}
	}

	problems := api.Group("/problems")
	problems.Use(middleware.JWTAuth())
	{
		problems.POST("", middleware.RequireRole("admin", "teacher"), controllers.CreateProblem)
		problems.GET("", controllers.GetProblems)
		problems.GET("/:id", controllers.GetProblemByID)
	}

	testcases := api.Group("/testcases")
	testcases.Use(middleware.JWTAuth())
	{
		testcases.POST("", middleware.RequireRole("admin", "teacher"), controllers.CreateTestcase)
	}

	submit := api.Group("/submit")
	submit.Use(middleware.JWTAuth())
	{
		submit.POST("", controllers.Submit)
	}
}
