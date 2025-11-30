package services

import (
	auth "confam-api/internal/auth"
	repositories "confam-api/internal/repositories"
	structs "confam-api/internal/structs"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IAuthService interface {
	Register(ctx context.Context, req structs.RegisterRequest) error
	Login(ctx context.Context, req structs.LoginRequest) (string, error)
	ForgotPassword(ctx context.Context, req structs.ForgotPasswordRequest) error
	PasswordReset(ctx context.Context, req structs.PasswordResetRequest) error
	ChangePassword(ctx context.Context, req structs.ChangePasswordRequest) error
}

// UserService implements the IUserService interface.
type AuthService struct {
	companyRepo repositories.ICompanyRepository
}

// NewUserService creates a new instance of UserService.
func NewAuthService(companyRepo repositories.ICompanyRepository) *AuthService {
	return &AuthService{companyRepo: companyRepo}
}

// Register contains the business logic for user registration.
func (s *AuthService) Register(ctx context.Context, req structs.RegisterRequest) error {
	// Business logic: check for existing company, hash password, create company record
	// You would call s.companyRepo.Create here
	// The controller is only responsible for receiving and sending data,
	// so the business logic should be here.
	return nil
}

// Login contains the business logic for user login and JWT generation.
func (s *AuthService) Login(ctx context.Context, req structs.LoginRequest) (string, error) {
	company, err := s.companyRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("Oops!, your email or password is invalid")
		}
		return "", errors.New("An internal server error occurred. Please try again later")
	}

	if err := auth.ComparePasswordAndHash(req.Password, *company.Password); err != nil {
		return "", errors.New("Oops!, your email or password is invalid")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": company.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New("Could not complete login. Please try again later.")
	}

	return tokenString, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req structs.ForgotPasswordRequest) error {
	// Business logic: check for existing company, hash password, create company record
	// You would call s.companyRepo.Create here
	// The controller is only responsible for receiving and sending data,
	// so the business logic should be here.
	return nil
}

func (s *AuthService) PasswordReset(ctx context.Context, req structs.PasswordResetRequest) error {
	// Business logic: check for existing company, hash password, create company record
	// You would call s.companyRepo.Create here
	// The controller is only responsible for receiving and sending data,
	// so the business logic should be here.
	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req structs.ChangePasswordRequest) error {
	// Business logic: check for existing company, hash password, create company record
	// You would call s.companyRepo.Create here
	// The controller is only responsible for receiving and sending data,
	// so the business logic should be here.
	return nil
}
