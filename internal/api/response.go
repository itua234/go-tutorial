package api

import (
	structs "confam-api/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

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
