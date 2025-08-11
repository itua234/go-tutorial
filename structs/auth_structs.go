package structs

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
