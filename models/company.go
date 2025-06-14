package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Company struct {
	ID                   string     `gorm:"type:char(36);primaryKey" json:"id"`
	Name                 string     `gorm:"type:varchar(255);not null" json:"name"`
	Logo                 *string    `gorm:"type:varchar(255)" json:"logo,omitempty"`
	Email                string     `gorm:"type:varchar(191);unique;not null" json:"email"`
	Password             *string    `gorm:"type:varchar(255)" json:"-"`
	Domain               *string    `gorm:"type:varchar(255)" json:"domain,omitempty"`
	WebhookURL           *string    `gorm:"type:varchar(255)" json:"webhook_url,omitempty"`
	Verified             bool       `gorm:"default:false" json:"verified"`
	NotificationsEnabled bool       `gorm:"default:true" json:"notifications_enabled"`
	EmailVerifiedAt      *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt            time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}

// BeforeCreate GORM hook to hash password before creating
func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	if company.ID == "" {
		company.ID = uuid.New().String()
	}
	if company.Password != nil && *company.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*company.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		pw := string(hashed)
		company.Password = &pw
	}
	if company.Email != "" {
		company.Email = lower(company.Email)
	}
	return nil
}

func (company *Company) BeforeUpdate(tx *gorm.DB) (err error) {
	if company.Password != nil && *company.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*company.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		pw := string(hashed)
		company.Password = &pw
	}
	if company.Email != "" {
		company.Email = lower(company.Email)
	}
	return nil
}

// Helper to lowercase email
func lower(s string) string {
	return string([]rune(s))
	// Or simply: return strings.ToLower(s)
}
