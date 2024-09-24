package middleware

import (
	"net/http"
	"pcgamedb/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Logger.Error("Recovery", zap.Any("error", rec), zap.Stack("stacktrace"))
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
			}
		}()
		c.Next()
	}
}
