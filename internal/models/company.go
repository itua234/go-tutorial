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
	Logo                 *string    `gorm:"type:varchar(255)" json:"logo"`
	Email                string     `gorm:"type:varchar(191);unique;not null" json:"email"`
	Password             *string    `gorm:"type:varchar(255)" json:"password"`
	Domain               *string    `gorm:"type:varchar(255)" json:"domain"`
	WebhookURL           *string    `gorm:"type:varchar(255)" json:"webhook_url"`
	Verified             bool       `gorm:"default:false" json:"verified"`
	NotificationsEnabled bool       `gorm:"default:true" json:"notifications_enabled"`
	EmailVerifiedAt      *time.Time `json:"email_verified_at"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}

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
