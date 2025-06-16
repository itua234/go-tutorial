package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"confam-api/database"
	"confam-api/models"

	"confam-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// import your encryption and token helpers
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
type KYCRequestInput struct {
	Customer     CustomerInput `json:"customer"`
	Reference    string        `json:"reference"`
	RedirectURL  string        `json:"redirect_url"`
	KYCLevel     string        `json:"kyc_level"`
	BankAccounts bool          `json:"bank_accounts"`
}

func InitiateKyc(c *gin.Context) {
	var input KYCRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": true})
		return
	}
	if input.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid reference, please retry with a unique reference", "error": true})
		return
	}
	cust := input.Customer
	if cust.Name == "" || cust.Email == "" || cust.Address == "" || cust.Identity.Type == "" || cust.Identity.Number == "" || input.RedirectURL == "" || input.KYCLevel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required customer or request fields.", "error": true})
		return
	}

	validIdentityTypes := []string{"BVN", "NIN"}
	found := false
	for _, t := range validIdentityTypes {
		if strings.ToUpper(cust.Identity.Type) == t {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid identity type: " + cust.Identity.Type + ". Must be one of BVN, NIN.",
			"error":   true,
		})
		return
	}

	// Generate token and encrypt customer data
	kycToken := utils.GenerateToken() // implement this
	customerBytes, _ := json.Marshal(cust)
	encryptedData := utils.Encrypt(string(customerBytes))

	// Get app from context (set by middleware)
	appVal, exists := c.Get("app")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "App not found in context", "error": true})
		return
	}
	app := appVal.(models.App)
	webhookURL := ""
	if app.WebhookURL != nil {
		webhookURL = *app.WebhookURL
	}

	// Lookup or create customer
	var customer models.Customer
	hashed := utils.HashFunction(cust.Email)
	emailHash := ""
	if hashed != nil {
		emailHash = *hashed
	}
	if err := database.DB.Where("email_hash = ?", emailHash).First(&customer).Error; err != nil {
		// Not found, create new
		customer = models.Customer{
			Token:            utils.GenerateToken(),
			Email:            cust.Email,
			EmailHash:        emailHash,
			Status:           "pending",
			KYCLevelAchieved: "none",
			IsBlacklisted:    false,
		}
		if err := database.DB.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create customer", "error": true})
			return
		}
	}

	// Create request in transaction
	var request models.Request
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		request = models.Request{
			CompanyID:             app.CompanyID,
			Reference:             input.Reference,
			RedirectURL:           input.RedirectURL,
			KYCLevel:              input.KYCLevel,
			BankAccountsRequested: input.BankAccounts,
			KYCToken:              kycToken,
			TokenExpiresAt:        time.Now().Add(1 * time.Hour),
			EncryptedData:         &encryptedData,
		}
		if err := tx.Create(&request).Error; err != nil {
			return err
		}
		// Optionally store token mapping in Redis or elsewhere
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request", "error": true})
		return
	}

	// Prepare response data
	data := gin.H{
		"id":             request.KYCToken,
		"customer":       customer.Token, // or request.KYCToken if no token
		"allow_url":      request.AllowURL,
		"reference":      request.Reference,
		"redirect_url":   request.RedirectURL,
		"bank_accounts":  request.BankAccountsRequested,
		"kyc_level":      request.KYCLevel,
		"is_blacklisted": false,
	}

	// Send webhook (non-blocking, just log for now)
	if webhookURL != "" {
		go func() {
			payload := gin.H{
				"event": "kyc.initiation.requested",
				"data": gin.H{
					"app":            app.ID,
					"business":       app.CompanyID,
					"id":             request.ID,
					"status":         request.Status,
					"reference":      request.Reference,
					"created_at":     request.CreatedAt,
					"kyc_level":      request.KYCLevel,
					"bank_accounts":  request.BankAccountsRequested,
					"is_blacklisted": false,
					"meta":           gin.H{},
				},
			}
			log.Println(payload)
			// Implement webhook sending logic here (e.g., http.Post)
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC process initiated successfully",
		"results": data,
		"error":   false,
	})
}
