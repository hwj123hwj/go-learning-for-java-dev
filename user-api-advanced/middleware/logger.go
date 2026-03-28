package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// 只在有错误时才附加错误信息
		if errMsg := c.Errors.String(); errMsg != "" {
			log.Printf("[%s] %s %s %d %v | %s", method, path, clientIP, statusCode, latency, errMsg)
		} else {
			log.Printf("[%s] %s %s %d %v", method, path, clientIP, statusCode, latency)
		}
	}
}
