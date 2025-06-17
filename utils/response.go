// pkg/utils/response.go
package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Message string      `json:"message"`
	Results interface{} `json:"results,omitempty"`
	Error   bool        `json:"error"`
}

// RespondWithError sends an error JSON response
func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Message: message,
		Error:   true,
	})
}

// RespondWithJSON sends a success JSON response
func RespondWithJSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
