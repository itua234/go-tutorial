package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// A very basic in-memory rate limiter example
var (
	requestCounts = make(map[string]int)
	lastReset     = time.Now()
	mu            sync.Mutex
	maxRequests   = 5 // 5 requests per 10 seconds
	resetInterval = 10 * time.Second
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		// Reset counts periodically
		if time.Since(lastReset) > resetInterval {
			requestCounts = make(map[string]int)
			lastReset = time.Now()
		}

		requestCounts[ip]++
		if requestCounts[ip] > maxRequests {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
