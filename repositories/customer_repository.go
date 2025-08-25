package repositories

import (
	"confam-api/models"
	"context"

	"gorm.io/gorm"
)

type ICustomerRepository interface {
	FindByID(ctx context.Context, id string) (*models.Customer, error)
	FindByEmail(ctx context.Context, email string) (*models.Customer, error)
	FindByEmailHash(email_hash string) (*models.Customer, error)
	Create(customer *models.Customer) error
	CreateIdentity(identity *models.Identity) error
}

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// FindByEmail finds a Customer by its email.
func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	result := r.db.WithContext(ctx).First(&customer, "email = ?", email)
	return &customer, result.Error
}

// FindByEmailHash finds a Customer by its email_hash.
func (r *CustomerRepository) FindByEmailHash(email_hash string) (*models.Customer, error) {
	var customer models.Customer
	result := r.db.First(&customer, "email_hash = ?", email_hash)
	return &customer, result.Error
}

// Create saves a new Customer record to the database.
func (r *CustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) CreateIdentity(identity *models.Identity) error {
	//return r.db.Create(identity).Error
	return r.db.Create(identity).Error
}
