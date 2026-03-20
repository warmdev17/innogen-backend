package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/controllers"
	"github.com/warmdev17/innogen-backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Authentication routes
		// NOTE: Register and send-otp routes are temporarily disabled
		// Only login is available for existing accounts
		api.POST("/auth/login", controllers.Login)
		api.POST("/auth/refresh", controllers.RefreshToken)

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Public Problems
		api.GET("/problems", controllers.GetProblems)
		api.GET("/problems/:id", controllers.GetProblemByID)

		// Course Structure
		api.GET("/subjects", controllers.GetSubjects)
		api.GET("/subjects/:slug", controllers.GetSubject)
		api.GET("/sessions/:id", controllers.GetSession)
		api.GET("/lessons/:id", controllers.GetLesson)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// User info
			me := protected.Group("/me")
			{
				me.GET("", controllers.GetCurrentUser)
			}

			// Auth
			auth := protected.Group("/auth")
			{
				auth.POST("/logout", controllers.Logout)
				auth.POST("/logout-all", controllers.LogoutAll)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin", "teacher"))
			{
				admin.POST("/problems", controllers.CreateProblem)
			}

			// Submissions
			submit := protected.Group("/submit")
			{
				submit.POST("", controllers.Submit)
				submit.GET("/:id", controllers.GetSubmission)
			}

			// Run code directly
			run := protected.Group("/run")
			{
				run.POST("", controllers.RunCode)
			}
		}
	}
}
