package handlers

import (
	"net/http"
	
	"your-app/internal/api"      // For sending back the error response
	"your-app/internal/validate" // The new validation package
)

// Example struct representing the incoming request body
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
	Email    string `json:"email" validate:"required,email"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	// 1. Decode Request Body (omitted for brevity)
	// ...

	// 2. Run Validation
	if errors := validate.Struct(req); errors != nil {
		// Use the new API error helper (Status 400 Bad Request)
		api.Error(w, http.StatusBadRequest, "Validation failed", errors) 
		return
	}
	
	// 3. Run Complex Business Rules (optional)
	if !validate.CheckUniqueUsername(req.Username) {
	    api.Error(w, http.StatusConflict, "Username is already taken")
	    return
	}

	// If everything passed, continue with service logic...
	// service.CreateUser(r.Context(), req)
}