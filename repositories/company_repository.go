package repositories

import (
	"confam-api/models"
	"context"

	"gorm.io/gorm"
)

type ICompanyRepository interface {
	FindByID(ctx context.Context, id string) (*models.Company, error)
	FindByEmail(ctx context.Context, email string) (*models.Company, error)
	Create(ctx context.Context, company *models.Company) error
}

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) FindByID(ctx context.Context, id string) (*models.Company, error) {
	var company models.Company
	if err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// FindByEmail finds a company by its email.
func (r *CompanyRepository) FindByEmail(ctx context.Context, email string) (*models.Company, error) {
	var company models.Company
	result := r.db.WithContext(ctx).First(&company, "email = ?", email)
	return &company, result.Error
}

// Create saves a new company record to the database.
func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}
