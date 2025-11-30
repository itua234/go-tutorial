package repositories

import (
	"confam-api/internal/models"
	"context"

	"gorm.io/gorm"
)

type ICompanyRepository interface {
	Create(ctx context.Context, company *models.Company) error
	FindByID(ctx context.Context, id string) (*models.Company, error)
	FindByEmail(ctx context.Context, email string) (*models.Company, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Company, error)
	Update(ctx context.Context, company *models.Company) error
	Delete(ctx context.Context, id string) error
	// Utility
	Count(ctx context.Context) (int64, error)
}

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create saves a new company record to the database.
func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
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

func (r *CompanyRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).
		Save(company).Error
}

func (r *CompanyRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Company{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *CompanyRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Delete(&models.Company{}, "id = ?", id).Error
}

func (r *CompanyRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Company{}).
		Count(&count).Error
	return count, err
}
