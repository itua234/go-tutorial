package repositories

import (
	"confam-api/models"

	"gorm.io/gorm"
)

type IRequestRepository interface {
	FindByToken(kyc_token string) (*models.Request, error)
	CountByReference(reference string) (int64, error)
	FindByReference(reference string) (*models.Request, error)
	Create(request *models.Request) error
}

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) FindByToken(kyc_token string) (*models.Request, error) {
	var request models.Request
	result := r.db.First(&request, "kyc_token = ?", kyc_token)
	return &request, result.Error
}

func (r *RequestRepository) CountByReference(reference string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Request{}).
		Where("reference = ?", reference).
		Count(&count).Error
	return count, err
}

func (r *RequestRepository) FindByReference(reference string) (*models.Request, error) {
	var request models.Request
	result := r.db.First(&request, "reference = ?", reference)
	return &request, result.Error
}

func (r *RequestRepository) Create(request *models.Request) error {
	return r.db.Create(request).Error
}
