package validate

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Global validator instance for reuse.
var v = validator.New()

func Struct(s interface{}) map[string]string {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string]string)

	// Iterate through the errors and format them
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		// Example message: "email is not a valid email address"
		validationErrors[fieldName] = fmt.Sprintf("%s failed validation rule: %s", fieldName, err.Tag())
	}

	return validationErrors
}

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
		case "string":
			errors[field] = "The " + readable + " attribute must be a string."
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
