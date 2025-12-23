package middleware

import (
	"dormitory-management/config"
	"dormitory-management/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := utils.ValidateToken(tokenString, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Akses ditolak. Hanya untuk administrator.",
				"redirect": "/student",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func StudentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "student" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Akses ditolak. Hanya untuk mahasiswa.",
				"redirect": "/admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
