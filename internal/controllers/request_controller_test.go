package controllers

// import (
//     "errors"
//     "testing"
//     "confam-api/models"
//     "confam-api/repositories/mocks"
//     "gorm.io/gorm"
// )

// func TestGetRequestByToken_Success(t *testing.T) {
//     // Arrange: Create mock that returns a valid request
//     expectedRequest := &models.Request{
//         ID:        "test-id-123",
//         KYCToken:  "valid-token",
//         Reference: "REF-123",
//         Status:    models.RequestStatusInitiated,
//     }

//     mockRepo := &mocks.MockRequestRepository{
//         FindByTokenFunc: func(kyc_token string) (*models.Request, error) {
//             return expectedRequest, nil
//         },
//     }

//     controller := NewRequestController(mockRepo)

//     // Act
//     result, err := controller.GetRequestByToken("valid-token")

//     // Assert
//     if err != nil {
//         t.Fatalf("Expected no error, got %v", err)
//     }

//     if result == nil {
//         t.Fatal("Expected result, got nil")
//     }

//     if result.ID != expectedRequest.ID {
//         t.Errorf("Expected ID %v, got %v", expectedRequest.ID, result.ID)
//     }

//     // Verify mock was called correctly
//     if !mockRepo.FindByTokenCalled {
//         t.Error("Expected FindByToken to be called")
//     }

//     if mockRepo.FindByTokenCalledWith != "valid-token" {
//         t.Errorf("Expected FindByToken called with 'valid-token', got %v",
//             mockRepo.FindByTokenCalledWith)
//     }
// }

// func TestGetRequestByToken_EmptyToken(t *testing.T) {
//     mockRepo := &mocks.MockRequestRepository{}
//     controller := NewRequestController(mockRepo)

//     result, err := controller.GetRequestByToken("")

//     if err == nil {
//         t.Fatal("Expected error for empty token, got nil")
//     }

//     if result != nil {
//         t.Error("Expected nil result for empty token")
//     }

//     // Should not call repository
//     if mockRepo.FindByTokenCalled {
//         t.Error("FindByToken should not be called with empty token")
//     }
// }

// func TestGetRequestByToken_NotFound(t *testing.T) {
//     mockRepo := &mocks.MockRequestRepository{
//         FindByTokenFunc: func(kyc_token string) (*models.Request, error) {
//             return nil, gorm.ErrRecordNotFound
//         },
//     }

//     controller := NewRequestController(mockRepo)

//     result, err := controller.GetRequestByToken("nonexistent-token")

//     if err == nil {
//         t.Fatal("Expected error when request not found")
//     }

//     if result != nil {
//         t.Error("Expected nil result when request not found")
//     }

//     expectedError := "request not found"
//     if err.Error() != expectedError {
//         t.Errorf("Expected error message '%v', got '%v'", expectedError, err.Error())
//     }
// }

// func TestGetRequestByToken_DatabaseError(t *testing.T) {
//     mockRepo := &mocks.MockRequestRepository{
//         FindByTokenFunc: func(kyc_token string) (*models.Request, error) {
//             return nil, errors.New("database connection failed")
//         },
//     }

//     controller := NewRequestController(mockRepo)

//     result, err := controller.GetRequestByToken("some-token")

//     if err == nil {
//         t.Fatal("Expected error when database fails")
//     }

//     if result != nil {
//         t.Error("Expected nil result when database fails")
//     }
// }

// func TestCheckReferenceExists_True(t *testing.T) {
//     mockRepo := &mocks.MockRequestRepository{
//         CountByReferenceFunc: func(reference string) (int64, error) {
//             return 1, nil
//         },
//     }

//     controller := NewRequestController(mockRepo)

//     exists, err := controller.CheckReferenceExists("REF-123")

//     if err != nil {
//         t.Fatalf("Expected no error, got %v", err)
//     }

//     if !exists {
//         t.Error("Expected reference to exist")
//     }

//     if !mockRepo.CountByReferenceCalled {
//         t.Error("Expected CountByReference to be called")
//     }
// }

// func TestCheckReferenceExists_False(t *testing.T) {
//     mockRepo := &mocks.MockRequestRepository{
//         CountByReferenceFunc: func(reference string) (int64, error) {
//             return 0, nil
//         },
//     }

//     controller := NewRequestController(mockRepo)

//     exists, err := controller.CheckReferenceExists("REF-999")

//     if err != nil {
//         t.Fatalf("Expected no error, got %v", err)
//     }

//     if exists {
//         t.Error("Expected reference to not exist")
//     }
// }
