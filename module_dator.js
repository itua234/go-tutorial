// validators/auth_validator.go
package validators

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AuthValidator struct {
	validator *validator.Validate
}

func NewAuthValidator() *AuthValidator {
	v := validator.New()
	
	// Register custom validations
	v.RegisterValidation("strong_password", validateStrongPassword)
	v.RegisterValidation("alpha_space", validateAlphaSpace)
	v.RegisterValidation("phone", validatePhone)
	
	return &AuthValidator{validator: v}
}

func (av *AuthValidator) ValidateStruct(s interface{}) error {
	return av.validator.Struct(s)
}

// Custom validation functions
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	
	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}

func validateAlphaSpace(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return regexp.MustCompile(`^\+?[1-9]\d{8,14}$`).MatchString(phone)
}

// models/auth_models.go
package models

type RegisterRequest struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=50" validate:"alpha_space"`
	LastName        string `json:"last_name" binding:"required,min=2,max=50" validate:"alpha_space"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=100" validate:"strong_password"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Phone           string `json:"phone" binding:"omitempty" validate:"phone"`
	Age             int    `json:"age" binding:"omitempty,min=13,max=120"`
	Terms           bool   `json:"terms" binding:"required" validate:"eq=true"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *RegisterRequest) Validate() map[string]string {
	errors := make(map[string]string)
	
	// Password confirmation check
	if r.Password != r.ConfirmPassword {
		errors["confirm_password"] = "Password confirmation does not match"
	}
	
	// Terms acceptance check
	if !r.Terms {
		errors["terms"] = "You must accept the terms and conditions"
	}
	
	return errors
}

// utils/validation_utils.go
package utils

import (
	"strings"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			errors[field] = getErrorMessage(e)
		}
	}
	
	return errors
}

func getErrorMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	
	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return "Please provide a valid email address"
	case "min":
		if e.Kind().String() == "string" {
			return field + " must be at least " + e.Param() + " characters long"
		}
		return field + " must be at least " + e.Param()
	case "max":
		if e.Kind().String() == "string" {
			return field + " must not exceed " + e.Param() + " characters"
		}
		return field + " must not exceed " + e.Param()
	case "alpha_space":
		return field + " must contain only letters and spaces"
	case "phone":
		return "Please provide a valid phone number"
	case "strong_password":
		return "Password must contain at least 8 characters with uppercase, lowercase, number and special character"
	case "eq":
		if e.Field() == "Terms" {
			return "You must accept the terms and conditions"
		}
		return field + " must equal " + e.Param()
	default:
		return field + " is invalid"
	}
}

// handlers/auth_handler.go
package handlers

import (
	"net/http"
	"your-app/models"
	"your-app/utils"
	"your-app/validators"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	validator *validators.AuthValidator
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		validator: validators.NewAuthValidator(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	
	// Bind JSON to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "Invalid input data",
			"details": utils.FormatValidationErrors(err),
		})
		return
	}
	
	// Custom struct validation
	if err := h.validator.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "Invalid input data",
			"details": utils.FormatValidationErrors(err),
		})
		return
	}
	
	// Model-specific validation
	if customErrors := req.Validate(); len(customErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "Registration validation failed",
			"details": customErrors,
		})
		return
	}
	
	// Business logic validation (e.g., check if email exists)
	if exists, err := h.checkEmailExists(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Server Error",
			"message": "Unable to validate email",
		})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Validation Error",
			"message": "Email is already registered",
			"details": map[string]string{
				"email": "This email address is already in use",
			},
		})
		return
	}
	
	// Process registration
	user, err := h.processRegistration(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Server Error",
			"message": "Registration failed",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "Invalid input data",
			"details": utils.FormatValidationErrors(err),
		})
		return
	}
	
	// Authenticate user
	user, err := h.authenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication Error",
			"message": "Invalid email or password",
		})
		return
	}
	
	// Generate token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Server Error",
			"message": "Failed to generate authentication token",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// Simulate business logic methods
func (h *AuthHandler) checkEmailExists(email string) (bool, error) {
	// Database check logic here
	return false, nil
}

func (h *AuthHandler) processRegistration(req models.RegisterRequest) (interface{}, error) {
	// Registration processing logic
	return gin.H{"id": 1, "email": req.Email}, nil
}

func (h *AuthHandler) authenticateUser(email, password string) (interface{}, error) {
	// Authentication logic
	return gin.H{"id": 1, "email": email}, nil
}

func (h *AuthHandler) generateToken(userID interface{}) (string, error) {
	// JWT token generation
	return "jwt_token_here", nil
}


import (
	"your-app/handlers"
	"github.com/gin-gonic/gin"
)

// func main() {
// 	r := gin.Default()
	
// 	authHandler := handlers.NewAuthHandler()
	
// 	api := r.Group("/api/v1")
// 	{
// 		auth := api.Group("/auth")
// 		{
// 			auth.POST("/register", authHandler.Register)
// 			auth.POST("/login", authHandler.Login)
// 		}
// 	}
	
// 	r.Run(":8080")
// }