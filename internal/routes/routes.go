package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/controllers"
	"github.com/warmdev17/innogen-backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/auth/login", controllers.Login)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
        
        // Public Problems
        api.GET("/problems", controllers.GetProblems)
        api.GET("/problems/:id", controllers.GetProblemByID)

		protected := api.Group("/me")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("", controllers.GetCurrentUser)
		}
        
        // Admin routes
        admin := api.Group("/admin")
        admin.Use(middleware.JWTAuth(), middleware.RequireRole("admin", "teacher"))
        {
            admin.POST("/problems", controllers.CreateProblem)
        }

        // Submissions
        submit := api.Group("/submit")
        submit.Use(middleware.JWTAuth())
        {
            submit.POST("", controllers.Submit)
            submit.GET("/:id", controllers.GetSubmission)
        }
        
        // Run code directly
        run := api.Group("/run")
        run.Use(middleware.JWTAuth())
        {
            run.POST("", controllers.RunCode)
        }
	}
}
