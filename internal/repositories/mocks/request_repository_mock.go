package mocks

import "confam-api/internal/models"

type MockRequestRepository struct {
	FindByTokenFunc      func(kyc_token string) (*models.Request, error)
	CountByReferenceFunc func(reference string) (int64, error)
	FindByReferenceFunc  func(reference string) (*models.Request, error)
	CreateFunc           func(request *models.Request) error

	// Fields to track calls and arguments
	FindByTokenCalled     bool
	FindByTokenCalledWith string

	CountByReferenceCalled     bool
	CountByReferenceCalledWith string

	FindByReferenceCalled     bool
	FindByReferenceCalledWith string

	CreateCalled     bool
	CreateCalledWith *models.Request
}

func (m *MockRequestRepository) FindByToken(kyc_token string) (*models.Request, error) {
	m.FindByTokenCalled = true
	m.FindByTokenCalledWith = kyc_token

	if m.FindByTokenFunc != nil {
		return m.FindByTokenFunc(kyc_token)
	}
	return nil, nil
}

func (m *MockRequestRepository) CountByReference(reference string) (int64, error) {
	m.CountByReferenceCalled = true
	m.CountByReferenceCalledWith = reference

	if m.CountByReferenceFunc != nil {
		return m.CountByReferenceFunc(reference)
	}
	return 0, nil
}

func (m *MockRequestRepository) FindByReference(reference string) (*models.Request, error) {
	m.FindByReferenceCalled = true
	m.FindByReferenceCalledWith = reference

	if m.FindByReferenceFunc != nil {
		return m.FindByReferenceFunc(reference)
	}
	return nil, nil
}

func (m *MockRequestRepository) Create(request *models.Request) error {
	m.CreateCalled = true
	m.CreateCalledWith = request

	if m.CreateFunc != nil {
		return m.CreateFunc(request)
	}
	return nil
}
