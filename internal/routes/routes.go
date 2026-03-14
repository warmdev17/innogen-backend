// Package routes
package routes

import (
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/innogen-backend/internal/controllers"
	"github.com/warmdev17/innogen-backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	pistonURL := os.Getenv("PISTON_URL")
	if pistonURL == "" {
		pistonURL = "http://localhost:2000"
	}
	remote, _ := url.Parse(pistonURL)
	pistonProxy := httputil.NewSingleHostReverseProxy(remote)
	r.Any("/piston/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		pistonProxy.ServeHTTP(c.Writer, c.Request)
	})

	api := r.Group("/api")
	{
		// Tạm bỏ đăng ký với OTP
		// api.POST("/auth/send-otp", controllers.SendOTP)
		// api.POST("/auth/register", controllers.Register)
		api.POST("/auth/login", controllers.Login)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

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
