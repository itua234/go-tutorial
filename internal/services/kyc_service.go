package services

import (
	crypto "confam-api/internal/crypto"
	"confam-api/internal/models"
	repositories "confam-api/internal/repositories"
	structs "confam-api/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Identity struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}
type CustomerInput struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Address  string   `json:"address"`
	Identity Identity `json:"identity"`
}
type KycRequestInput struct {
	Customer     CustomerInput `json:"customer"`
	Reference    string        `json:"reference"`
	RedirectURL  string        `json:"redirect_url"`
	KYCLevel     string        `json:"kyc_level"`
	BankAccounts bool          `json:"bank_accounts"`
}

type IKycService interface {
	ValidateIdentityType(ctx context.Context, identityType string) bool
	FindOrCreateCustomer(ctx context.Context, req structs.CustomerInput) (*models.Customer, error)
	CreateKYCRequest(ctx context.Context, app models.App, req structs.KycRequestInput) (*models.Request, error)
	FetchKycRequest(ctx context.Context, kycToken string) (*models.Request, *models.Customer, error)
	IsReferenceUnique(ctx context.Context, reference string) (bool, error)
}

type KYCService struct {
	customerRepo repositories.ICustomerRepository
	requestRepo  repositories.IRequestRepository
}

func NewKYCService(
	customerRepo repositories.ICustomerRepository,
	requestRepo repositories.IRequestRepository,
) *KYCService {
	return &KYCService{
		customerRepo: customerRepo,
		requestRepo:  requestRepo,
	}
}

func (s *KYCService) ValidateIdentityType(ctx context.Context, identityType string) bool {
	validTypes := []string{"BVN", "NIN"}
	for _, t := range validTypes {
		if strings.EqualFold(t, identityType) {
			return true
		}
	}
	return false
}

func (s *KYCService) FindOrCreateCustomer(ctx context.Context, req structs.CustomerInput) (*models.Customer, error) {
	hash := crypto.HashSHA256(req.Email)
	customer, err := s.customerRepo.FindByEmailHash(hash)
	if err != nil {
		return nil, err
	}
	if customer != nil {
		return customer, nil
	}

	token, err := crypto.GenerateHexToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate customer token: %w", err)
	}
	customer = &models.Customer{
		Token:            token,
		Email:            req.Email,
		EmailHash:        hash,
		Status:           "pending",
		KYCLevelAchieved: "none",
		IsBlacklisted:    false,
	}
	if err := s.customerRepo.Create(customer); err != nil {
		return nil, err
	}

	encryptedValue, err := crypto.Encrypt(req.Identity.Number)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt identity number: %w", err)
	}

	identity := &models.Identity{
		CustomerID: customer.ID,
		Type:       models.IdentityType(req.Identity.Type),
		Value:      encryptedValue,
	}
	if err := s.customerRepo.CreateIdentity(identity); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *KYCService) CreateKYCRequest(
	ctx context.Context,
	app models.App,
	req structs.KycRequestInput,
) (*models.Request, error) {
	// Ensure unique reference
	// count, err := s.requestRepo.CountByReference(req.Reference)
	// if err != nil {
	// 	return nil, err
	// }
	// if count > 0 {
	// 	return nil, errors.New("duplicate reference")
	// }

	customerBytes, _ := json.Marshal(req.Customer)
	encryptedData, _ := crypto.Encrypt(string(customerBytes))

	kycToken, err := crypto.GenerateHexToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate KYC token: %w", err)
	}

	request := &models.Request{
		CompanyID:             app.CompanyID,
		Reference:             req.Reference,
		RedirectURL:           req.RedirectURL,
		KYCLevel:              req.KYCLevel,
		BankAccountsRequested: req.BankAccounts,
		KYCToken:              kycToken,
		TokenExpiresAt:        time.Now().Add(1 * time.Hour),
		EncryptedData:         &encryptedData,
	}

	if err := s.requestRepo.Create(request); err != nil {
		return nil, err
	}
	return request, nil
}

func (s *KYCService) FetchKycRequest(ctx context.Context, kycToken string) (*models.Request, *models.Customer, error) {
	// Use the repository to find the request by its token.
	request, err := s.requestRepo.FindByToken(kycToken)
	if err != nil {
		return nil, nil, err
	}

	// Business logic for status checks.
	if request.Status == "completed" {
		return nil, nil, fmt.Errorf("KYC request already completed")
	}
	if request.Status == "failed" {
		return nil, nil, fmt.Errorf("KYC request failed")
	}

	// Use the repository to find the associated customer.
	var decrypted map[string]interface{}
	if err := json.Unmarshal([]byte(*request.EncryptedData), &decrypted); err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt customer data")
	}

	email, ok := decrypted["email"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("customer email not found in decrypted data")
	}

	customer, err := s.customerRepo.FindByEmailHash(crypto.HashSHA256(email))
	if err != nil {
		return nil, nil, err
	}

	return request, customer, nil
}

func (s *KYCService) IsReferenceUnique(ctx context.Context, reference string) (bool, error) {
	// Use the repository method here.
	count, err := s.requestRepo.CountByReference(reference)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
