package controllers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"confam-api/models"
	"confam-api/structs"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestContext(method, url string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

// MockKycService is a mock implementation of the IKycService interface.
type MockKycService struct {
	mock.Mock
}

func (m *MockKycService) IsReferenceUnique(ctx context.Context, ref string) (bool, error) {
	args := m.Called(ctx, ref)
	return args.Bool(0), args.Error(1)
}

func (m *MockKycService) ValidateIdentityType(ctx context.Context, idType string) bool {
	args := m.Called(ctx, idType)
	return args.Bool(0)
}

func (m *MockKycService) FindOrCreateCustomer(ctx context.Context, customer structs.CustomerInput) (*models.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockKycService) CreateKYCRequest(ctx context.Context, app models.App, req structs.KycRequestInput) (*models.Request, error) {
	args := m.Called(ctx, app, req)
	return args.Get(0).(*models.Request), args.Error(1)
}

func (m *MockKycService) FetchKycRequest(ctx context.Context, kycToken string) (*models.Request, *models.Customer, error) {
	args := m.Called(ctx, kycToken)
	return args.Get(0).(*models.Request), args.Get(1).(*models.Customer), args.Error(2)
}

// MockWebhookService is a mock implementation of the IWebhookService interface.
type MockWebhookService struct {
	mock.Mock
}

func (m *MockWebhookService) SendWebhook(webhookURL, event string, data any) {
	m.Called(webhookURL, event, data)
}

func TestInitiateKyc_Success(t *testing.T) {
	mockKycService := new(MockKycService)
	mockWebhookService := new(MockWebhookService)
	ctrl := NewKycController(mockKycService, mockWebhookService)

	// Mock service dependencies
	mockKycService.On("IsReferenceUnique", mock.Anything, "unique-ref-123").Return(true, nil)
	mockKycService.On("ValidateIdentityType", mock.Anything, "BVN").Return(true)
	mockKycService.On("FindOrCreateCustomer", mock.Anything, mock.Anything).Return(&models.Customer{Token: "customer-token"}, nil)
	mockKycService.On("CreateKYCRequest", mock.Anything, mock.Anything, mock.Anything).Return(&models.Request{
		KYCToken:              "kyc-token-123",
		RedirectURL:           "http://redirect.url",
		Reference:             "unique-ref-123",
		KYCLevel:              "LEVEL_1",
		BankAccountsRequested: true,
	}, nil)
	mockWebhookService.On("SendWebhook", mock.Anything, mock.Anything, mock.Anything).Return()

	// Create a request body
	requestBody := []byte(`{
        "reference": "unique-ref-123",
        "redirect_url": "http://redirect.url",
        "kyc_level": "LEVEL_1",
        "customer": {
            "name": "John Doe",
            "email": "john.doe@example.com",
            "address": "123 Main St",
            "identity": {
                "type": "BVN",
                "number": "12345678901"
            }
        }
    }`)

	// Create a gin context
	c, w := createTestContext(http.MethodPost, "/kyc/initiate", requestBody)
	app := models.App{
		ID:         "app-id-1",
		CompanyID:  "company-id-1",
		WebhookURL: new(string),
	}
	*app.WebhookURL = "http://webhook.url"
	c.Set("app", app)

	// Call the handler
	ctrl.InitiateKyc(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"error":false,"message":"KYC process initiated successfully","results":{"allow_url":"","bank_accounts":true,"customer":"customer-token","id":"kyc-token-123","is_blacklisted":false,"kyc_level":"LEVEL_1","redirect_url":"http://redirect.url","reference":"unique-ref-123"}}`
	assert.JSONEq(t, expectedBody, w.Body.String())

	// Verify mock expectations
	mockKycService.AssertExpectations(t)
	mockWebhookService.AssertExpectations(t)
}
