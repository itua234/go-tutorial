package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

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

func SetupValidator() {
	validate = validator.New()

	// Register custom validation functions
	validate.RegisterValidation("alpha", validateAlpha)
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("strong_password", validateStrongPassword)
}
