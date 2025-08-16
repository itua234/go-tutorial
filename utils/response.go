package utils

import (
	structs "confam-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// // SuccessResponse sends a custom success response.
func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// // ErrorResponse sends a custom error response.
func ErrorResponse(c *gin.Context, status int, message string, errors interface{}) {
	c.JSON(status, structs.ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// // You can create more specific error helpers for common scenarios
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	})
}
