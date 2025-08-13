package utils

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ToSnakeCase converts a CamelCase string to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToHumanReadable converts a CamelCase string to a human-readable format with spaces.
func ToHumanReadable(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune(' ')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// FormatValidationErrors formats the validation errors into a map of field -> message.
func FormatValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, fieldError := range validationErrors {
		// field := strings.ToLower(fieldError.Field())
		// Use the Go field name and convert it to snake_case
		field := ToSnakeCase(fieldError.Field())
		// Use the Go field name and convert it to a human-readable format for the error message value
		readable := ToHumanReadable(fieldError.Field())
		switch fieldError.Tag() {
		case "required":
			errors[field] = readable + " field is required."
		case "email":
			errors[field] = readable + " must be a valid email address."
		case "min":
			errors[field] = readable + " must be at least " + fieldError.Param() + " characters"
		case "max":
			errors[field] = readable + " must not exceed " + fieldError.Param() + " characters"
		default:
			errors[field] = readable + " is invalid"
		}
	}
	return errors
}
