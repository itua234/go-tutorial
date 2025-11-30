package controllers

import (
	response "confam-api/internal/api"
	models "confam-api/internal/models"
	services "confam-api/internal/services"
	structs "confam-api/internal/structs"
	"confam-api/internal/validate"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

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
	var req structs.KycRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := validate.FormatValidationErrors(validationErrors)
			response.ValidationErrorResponse(c, errors)
			return
		}
		response.ErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	if req.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid reference, please retry with a unique reference", "error": true})
		return
	} else {
		isUnique, _ := ctrl.kycService.IsReferenceUnique(c, req.Reference)
		if !isUnique {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide a unique reference", "error": true})
			return
		}
	}
	cust := req.Customer
	if cust.Name == "" || cust.Email == "" || cust.Address == "" || cust.Identity.Type == "" || cust.Identity.Number == "" || req.RedirectURL == "" || req.KYCLevel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required customer or request fields.", "error": true})
		return
	}

	if !ctrl.kycService.ValidateIdentityType(c, req.Customer.Identity.Type) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid identity type: " + req.Customer.Identity.Type + ". Must be one of BVN, NIN.",
			"error":   true,
		})
		return
	}

	// Get app from context (set by middleware)
	appVal, _ := c.Get("app")
	app := appVal.(models.App)
	webhookURL := ""
	if app.WebhookURL != nil {
		webhookURL = *app.WebhookURL
	}

	// Lookup or create customer
	customer, err := ctrl.kycService.FindOrCreateCustomer(c, req.Customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process customer", "error": true})
		return
	}

	// Create request in transaction
	request, err := ctrl.kycService.CreateKYCRequest(c, app, req)
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

	//customer["bank_accounts"] = []string{}
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
