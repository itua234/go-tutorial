package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("--- INCOMING REQUEST ---\nOrigin: %s\nMethod: %s\nPath: %s\n",
			param.Request.Header.Get("Origin"),
			param.Method,
			param.Path,
		)
	})
}

func CorsHeadersLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log CORS headers after response
		fmt.Println("CORS Headers:")
		for key, value := range c.Writer.Header() {
			if strings.HasPrefix(key, "Access-Control-") {
				fmt.Printf("  %s: %s\n", key, strings.Join(value, ", "))
			}
		}
	}
}

// CORSMiddleware sets up CORS headers
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = c.GetHeader("origin")
		}

		//Get allowed origins from environment variable
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		var allowedOrigins []string
		if allowedOriginsEnv != "" {
			allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		}

		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.TrimSpace(allowedOrigin) == origin {
				isAllowed = true
				break
			}
		}

		// If the origin is NOT in the allowed list, abort the request.
		// This is the core logic to prevent all origins from being allowed.
		if !isAllowed && origin != "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: This origin is not allowed."})
			return
		}

		// If the origin is allowed, set the correct CORS headers
		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//c.Writer.Header().Set("Access-Control-Allow-Origin", origin)

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			fmt.Printf("OPTIONS request to: %s\n", c.Request.URL.Path)
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
