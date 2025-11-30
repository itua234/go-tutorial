package services

import (
	models "confam-api/internal/models"
	"context"
)

type ICompanyService interface {
	GetCompanyById(ctx context.Context, id string) (*models.Company, error)
	GetCompanyByEmail(ctx context.Context, email string) (*models.Company, error)
	CreateCompany(ctx context.Context, company *models.Company) error
}
