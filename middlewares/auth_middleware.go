package middlewares

import (
	client "confam-api/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"confam-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthenticateAppBySecretKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		rawToken := c.GetHeader("x-allow-key")
		token := ""
		if rawToken != "" {
			// If multiple headers, use the first
			if strings.Contains(rawToken, ",") {
				tokens := strings.Split(rawToken, ",")
				token = strings.TrimSpace(tokens[0])
				log.Printf("Multiple API keys found in 'x-allow-key' header, using the first one.")
			} else {
				token = rawToken
			}
		}
		requestMethod := c.Request.Method
		requestUrl := c.Request.RequestURI
		clientIp := c.ClientIP()
		log.Printf("Request received: %s %s from %s", requestMethod, requestUrl, clientIp)

		if token == "" {
			log.Printf("Authentication failed: Missing API key.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "unauthenticated"})
			c.Abort()
			return
		}

		log.Printf("Attempting authentication with key: %s... (masked)", maskToken(token))

		// Lookup appId in Redis
		appId, err := client.RedisClient.Get(c, "secret:"+token).Result()
		if err != nil || appId == "" {
			log.Printf("Authentication failed: Invalid API key - %s...", maskToken(token))
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid API key"})
			c.Abort()
			return
		}

		// Lookup App in DB, including Company
		var app models.App
		if err := db.Preload("Company").First(&app, "id = ?", appId).Error; err != nil {
			log.Printf("Authentication failed: App not found for ID: %s", appId)
			c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "App not found"})
			c.Abort()
			return
		}
		// Store app in context for downstream handlers
		c.Set("app", app)

		responseTime := time.Since(startTime).Milliseconds()
		companyName := "N/A"
		if app.Company != nil {
			companyName = app.Company.Name
		}
		log.Printf("Authentication successful for App ID: %s (Company: %s). Request processed in %dms.", appId, companyName, responseTime)
		c.Next()
	}
}

// Helper to mask token for logging
func maskToken(token string) string {
	if len(token) > 6 {
		return token[:6]
	}
	return token
}
