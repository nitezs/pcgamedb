package middleware

import (
	"pcgamedb/log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime).Milliseconds()
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		if strings.HasPrefix(reqURI, "/swagger/") ||
			strings.EqualFold(reqURI, "/favicon.ico") {
			return
		}

		log.Logger.Info(
			"request",
			zap.Int("code", statusCode),
			zap.String("method", reqMethod),
			zap.String("uri", reqURI),
			zap.String("ip", clientIP),
			zap.String("latency", strconv.Itoa(int(latencyTime))+"ms"),
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Logger.Error(e)
			}
		}
	}
}
