package controllers

import (
	"net/http"

	services "confam-api/services"
	structs "confam-api/structs"
	response "confam-api/utils"
	utils "confam-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthController handles HTTP requests.
type AuthController struct {
	AuthService services.IAuthService
}

// NewAuthController creates a new controller with a service dependency.
func NewAuthController(authService services.IAuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func Register(c *gin.Context) {
	var req structs.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := utils.FormatValidationErrors(validationErrors)
			response.ValidationErrorResponse(c, errors)
			return
		}
		response.ErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration Successful",
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req structs.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := utils.FormatValidationErrors(validationErrors)
			response.ValidationErrorResponse(c, errors)
			return
		}
		response.ErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// Pass the request data to the service and receive the token.
	token, err := ac.AuthService.Login(c, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Login Successful", gin.H{"token": token})
}

func ForgotPassword(c *gin.Context) {
	var req structs.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := utils.FormatValidationErrors(validationErrors)
			response.ValidationErrorResponse(c, errors)
			return
		}
		response.ErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// if err := ac.AuthService.ForgotPassword(c, req); err != nil {
	// 	response.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Forgot Password Successful",
	})
}

func PasswordReset(c *gin.Context) {
	var req structs.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := utils.FormatValidationErrors(validationErrors)
			response.ValidationErrorResponse(c, errors)
			return
		}
		response.ErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Forgot Password Successful",
	})
}

func ChangePassword(c *gin.Context) {
	var req structs.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := utils.FormatValidationErrors(validationErrors)
			response.ValidationErrorResponse(c, errors)
			return
		}
		response.ErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Change Password Successful",
	})
}
