// Package middleware
package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		if !slices.Contains(allowedRoles, userRole) {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden: You don't have permission"})
			return
		}

		c.Next()
	}
}
