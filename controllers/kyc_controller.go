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
	"bytes"
	"errors"
	"slices"
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

func InitiateKyc(c *gin.Context) {
	var input KycRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": true})
		return
	}

	if input.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid reference, please retry with a unique reference", "error": true})
		return
	} else {
		var count int64
		_ = database.DB.Model(&models.Request{}).Where("reference = ?", input.Reference).Count(&count).Error
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "please provide a unique reference", "error": true})
			return
		}
	}
	cust := input.Customer
	if cust.Name == "" || cust.Email == "" || cust.Address == "" || cust.Identity.Type == "" || cust.Identity.Number == "" || input.RedirectURL == "" || input.KYCLevel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required customer or request fields.", "error": true})
		return
	}

	validIdentityTypes := []string{"BVN", "NIN"}
	found := slices.Contains(validIdentityTypes, strings.ToUpper(cust.Identity.Type))
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid identity type: " + cust.Identity.Type + ". Must be one of BVN, NIN.",
			"error":   true,
		})
		return
	}

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

	//kycService := services.NewKYCService(database.DB)

	// Lookup or create customer
	customer, err := findOrCreateCustomer(input.Customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process customer", "error": true})
		return
	}

	// Create request in transaction
	request, err := createKYCRequest(input, app)
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
		data := map[string]interface{}{
			"app":            app.ID,
			"business":       app.CompanyID,
			"id":             request.ID,
			"status":         request.Status,
			"reference":      request.Reference,
			"created_at":     request.CreatedAt,
			"kyc_level":      request.KYCLevel,
			"bank_accounts":  request.BankAccountsRequested,
			"is_blacklisted": false,
			"meta":           map[string]interface{}{},
		}
		log.Println(data)
		//go sendWebhook(webhookURL, "kyc.initiation.requested", data) // Non-blocking
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC process initiated successfully",
		"results": data,
		"error":   false,
	})
}

func FetchKycRequest(c *gin.Context) {
	kyc_token := c.Param("kyc_token")
	if kyc_token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Missing KYC token",
		})
		return
	}

	var request models.Request
	result := database.DB.Where("kyc_token = ?", kyc_token).First(&request)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "KYC request not found",
		})
		return
	} else if request.Status == "completed" {
		log.Printf("KYC request already completed: %s", kyc_token)
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "KYC request already completed.",
		})
		return
	} else if request.Status == "failed" {
		log.Printf("KYC request failed: %s", kyc_token)
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "KYC request failed. Please try again.",
		})
		return
	}

	var decrypted map[string]interface{}
	json.Unmarshal([]byte(*request.EncryptedData), &decrypted)
	var customer models.Customer
	database.DB.Preload("Identities").
		Where("email_hash = ?", utils.HashFunction(decrypted["email"].(string))).
		First(&customer)

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC request fetched successfully",
		"results": gin.H{
			"redirect_url":            request.RedirectURL,
			"kyc_level":               request.KYCLevel,
			"bank_accounts_requested": request.BankAccountsRequested,
			"customer":                customer,
		},
		"error": false,
	})
}

func findOrCreateCustomer(input CustomerInput) (*models.Customer, error) {
	// ...lookup or create customer...
	var customer models.Customer
	emailHash := utils.HashFunction(input.Email)
	// Try to find customer by email_hash
	result := database.DB.Where("email_hash = ?", emailHash).First(&customer)
	if result.Error == nil {
		// Found existing customer
		return &customer, nil
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Some DB error occurred
		return nil, result.Error
	}

	// Not found, create new
	customer = models.Customer{
		Token:            utils.GenerateToken(),
		Email:            input.Email,
		EmailHash:        emailHash,
		Status:           "pending",
		KYCLevelAchieved: "none",
		IsBlacklisted:    false,
	}
	if err := database.DB.Create(&customer).Error; err != nil {
		return nil, err
	}

	identity := models.Identity{
		CustomerID: customer.ID,
		Type:       models.IdentityType(input.Identity.Type),
		Value:      utils.Encrypt(input.Identity.Number),
	}
	database.DB.Create(&identity)

	return &customer, nil
}

func createKYCRequest(
	req KycRequestInput,
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
	if err := database.DB.Create(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func sendWebhook(webhookURL, event string, data any) {
	payload := map[string]any{
		"event": event,
		"data":  data,
	}

	// Log the payload
	log.Println(payload)

	// Marshal to JSON
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal webhook payload:", err)
		return
	}

	// Send POST request
	response, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("Failed to send webhook:", err)
		return
	}
	defer response.Body.Close()
	log.Println("Webhook sent, status:", response.Status)
}
