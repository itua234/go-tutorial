package validator

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// User structs for validation
type RegisterRequest struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=50" validate:"alpha"`
	LastName        string `json:"last_name" binding:"required,min=2,max=50" validate:"alpha"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Phone           string `json:"phone" binding:"omitempty" validate:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Custom validator instance
var validate *validator.Validate

// Initialize custom validators
func init() {
	validate = validator.New()

	// Register custom validation functions
	validate.RegisterValidation("alpha", validateAlpha)
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// Custom validation functions
func validateAlpha(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Simple phone validation - adjust regex as needed
	return regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(phone)
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// At least 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special char
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}

// Error response structure
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Validation helper functions
func getValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must not exceed " + e.Param() + " characters"
			case "alpha":
				errors[field] = field + " must contain only letters and spaces"
			case "phone":
				errors[field] = "Invalid phone number format"
			case "strong_password":
				errors[field] = "Password must contain at least 8 characters with uppercase, lowercase, number and special character"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}

	return errors
}

// Custom validation for password confirmation
func validatePasswordMatch(req RegisterRequest) map[string]string {
	errors := make(map[string]string)

	if req.Password != req.ConfirmPassword {
		errors["confirm_password"] = "Password confirmation does not match"
	}

	return errors
}

// Middleware for handling validation errors
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			if err.Type == gin.ErrorTypeBind {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Validation Error",
					Message: "Invalid input data",
					Details: getValidationErrors(err.Err),
				})
				return
			}
		}
	}
}

// Registration handler
func registerHandler(c *gin.Context) {
	var req RegisterRequest

	// Bind JSON to struct and validate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation Error",
			Message: "Invalid input data",
			Details: getValidationErrors(err),
		})
		return
	}

	// Additional custom validations
	customErrors := make(map[string]string)

	// Check password confirmation
	if passwordErrors := validatePasswordMatch(req); len(passwordErrors) > 0 {
		for k, v := range passwordErrors {
			customErrors[k] = v
		}
	}

	// Custom validation using validator package
	if err := validate.Struct(req); err != nil {
		for k, v := range getValidationErrors(err) {
			customErrors[k] = v
		}
	}

	// Check for duplicate email (simulate database check)
	if isEmailTaken(req.Email) {
		customErrors["email"] = "Email is already registered"
	}

	// If there are custom validation errors, return them
	if len(customErrors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation Error",
			Message: "Registration failed",
			Details: customErrors,
		})
		return
	}

	// Process registration (save to database, hash password, etc.)
	user := processRegistration(req)

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Registration successful",
		Data:    user,
	})
}

// Login handler
func loginHandler(c *gin.Context) {
	var req LoginRequest

	// Bind and validate JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation Error",
			Message: "Invalid input data",
			Details: getValidationErrors(err),
		})
		return
	}

	// Additional business logic validation
	if !isValidCredentials(req.Email, req.Password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Authentication Error",
			Message: "Invalid email or password",
		})
		return
	}

	// Generate token and return success response
	token := generateAuthToken(req.Email)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Login successful",
		Data: gin.H{
			"token": token,
			"user":  getUserByEmail(req.Email),
		},
	})
}

// Simulate database/business logic functions
func isEmailTaken(email string) bool {
	// Simulate database check
	takenEmails := []string{"taken@example.com", "admin@example.com"}
	for _, taken := range takenEmails {
		if email == taken {
			return true
		}
	}
	return false
}

func isValidCredentials(email, password string) bool {
	// Simulate credential validation
	// In real app, hash password and compare with database
	return email == "user@example.com" && password == "password123"
}

func processRegistration(req RegisterRequest) interface{} {
	// Simulate user creation
	return gin.H{
		"id":         123,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      req.Email,
	}
}

func generateAuthToken(email string) string {
	// Simulate JWT token generation
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

func getUserByEmail(email string) interface{} {
	// Simulate user retrieval
	return gin.H{
		"id":         123,
		"first_name": "John",
		"last_name":  "Doe",
		"email":      email,
	}
}

// Rate limiting middleware (optional)
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement rate limiting logic here
		// For example, using redis or in-memory store
		c.Next()
	}
}

// Main function to set up routes
func main() {
	r := gin.Default()

	// Apply validation error handler middleware
	r.Use(ValidationErrorHandler())

	// Auth routes with rate limiting
	auth := r.Group("/api/auth")
	auth.Use(rateLimitMiddleware())
	{
		auth.POST("/register", registerHandler)
		auth.POST("/login", loginHandler)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.Run(":8080")
}
