package middleware

import (
	"net/http"
	"pcgamedb/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	apiKey := config.Config.Server.SecretKey
	if apiKey == "" {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "API key is not configured properly.",
			})
			c.Abort()
		}
	}
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized. No API key provided.",
			})
			c.Abort()
			return
		}
		if strings.TrimPrefix(auth, "Bearer ") != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized. Invalid API key.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
