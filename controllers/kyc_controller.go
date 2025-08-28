package controllers

import (
	"errors"
	"log"
	"net/http"

	"confam-api/database"
	models "confam-api/models"
	services "confam-api/services"
	structs "confam-api/structs"

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
type KycRequestInput struct {
	Customer     CustomerInput `json:"customer"`
	Reference    string        `json:"reference"`
	RedirectURL  string        `json:"redirect_url"`
	KYCLevel     string        `json:"kyc_level"`
	BankAccounts bool          `json:"bank_accounts"`
}

type kycController struct {
	kycService     services.IKycService
	webhookService services.IWebhookService
}

func NewKycController(
	kycService services.IKycService,
	webhookService services.IWebhookService,
) *kycController {
	return &kycController{
		kycService:     kycService,
		webhookService: webhookService,
	}
}

func (ctrl *kycController) InitiateKyc(c *gin.Context) {
	var input structs.KycRequestInput
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

	if !ctrl.kycService.ValidateIdentityType(c, input.Customer.Identity.Type) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid identity type: " + input.Customer.Identity.Type + ". Must be one of BVN, NIN.",
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

	// Lookup or create customer
	customer, err := ctrl.kycService.FindOrCreateCustomer(c, input.Customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process customer", "error": true})
		return
	}

	// Create request in transaction
	request, err := ctrl.kycService.CreateKYCRequest(c, app, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request", "error": true})
		return
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
		go ctrl.webhookService.SendWebhook(webhookURL, "kyc.initiation.requested", data)
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

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC process initiated successfully",
		"results": data,
		"error":   false,
	})
}

func (ctrl *kycController) FetchKycRequest(c *gin.Context) {
	kyc_token := c.Param("kyc_token")
	if kyc_token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Missing KYC token",
		})
		return
	}

	// Call the service to handle all the core logic.
	request, customer, err := ctrl.kycService.FetchKycRequest(c, kyc_token)
	if err != nil {
		log.Printf("Error fetching KYC request: %v", err)

		// Use a switch statement to handle different service-level errors.
		var status int
		var message string

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			status = http.StatusNotFound
			message = "KYC request not found"
		case err.Error() == "KYC request already completed":
			status = http.StatusForbidden
			message = "KYC request already completed."
		case err.Error() == "KYC request failed":
			status = http.StatusForbidden
			message = "KYC request failed. Please try again."
		default:
			status = http.StatusInternalServerError
			message = "Failed to process request"
		}

		c.JSON(status, gin.H{
			"error":   true,
			"message": message,
		})
		return
	}

	// Prepare the final response data based on the results from the service.
	c.JSON(http.StatusOK, gin.H{
		"message": "KYC request fetched successfully",
		"results": gin.H{
			"redirect_url":            request.RedirectURL,
			"kyc_level":               request.KYCLevel,
			"bank_accounts_requested": request.BankAccountsRequested,
			"customer":                customer,
			"status":                  request.Status,
		},
		"error": false,
	})
}
