package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Process request

		duration := time.Since(start)
		log.Printf("Method: %s | Path: %s | Status: %d | Latency: %s\n",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
