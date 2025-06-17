package services

import (
	"bytes"
	"confam-api/models"
	"confam-api/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"confam-api/dto"

	"gorm.io/gorm"
)

type KYCService struct {
	DB *gorm.DB
}

func NewKYCService(db *gorm.DB) *KYCService {
	return &KYCService{DB: db}
}

func (s *KYCService) FindOrCreateCustomer(
	input dto.CustomerInput,
) (*models.Customer, error) {
	var customer models.Customer
	emailHash := utils.HashFunction(input.Email)
	result := s.DB.Where("email_hash = ?", emailHash).First(&customer)
	if result.Error == nil {
		return &customer, nil
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	customer = models.Customer{
		Token:            utils.GenerateToken(),
		Email:            input.Email,
		EmailHash:        emailHash,
		Status:           "pending",
		KYCLevelAchieved: "none",
		IsBlacklisted:    false,
	}
	if err := s.DB.Create(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *KYCService) CreateKYCRequest(
	req dto.KycRequestInput,
	app models.App,
) (*models.Request, error) {
	var request models.Request
	kycToken := utils.GenerateToken()
	customerBytes, _ := json.Marshal(req.Customer)
	encryptedData := utils.Encrypt(string(customerBytes))
	request = models.Request{
		CompanyID:             app.CompanyID,
		Reference:             req.Reference,
		RedirectURL:           req.RedirectURL,
		KYCLevel:              req.KYCLevel,
		BankAccountsRequested: req.BankAccounts,
		KYCToken:              kycToken,
		TokenExpiresAt:        time.Now().Add(1 * time.Hour),
		EncryptedData:         &encryptedData,
	}
	if err := s.DB.Create(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (s *KYCService) SendWebhook(
	webhookURL,
	event string,
	data any,
) {
	payload := map[string]any{
		"event": event,
		"data":  data,
	}
	log.Println(payload)
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal webhook payload:", err)
		return
	}
	response, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("Failed to send webhook:", err)
		return
	}
	defer response.Body.Close()
	log.Println("Webhook sent, status:", response.Status)
}
