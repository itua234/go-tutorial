package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Decrypt(value string) string {
	return value
}

type Status string
type AccessType string

const (
	KYCNone  KYCLevel = "none"
	KYCTier1 KYCLevel = "tier_1"
	KYCTier2 KYCLevel = "tier_2"
	KYCTier3 KYCLevel = "tier_3"

	StatusPending  Status = "pending"
	StatusVerified Status = "verified"
	StatusRejected Status = "rejected"

	AccessTemporary AccessType = "temporary"
	AccessPermanent AccessType = "permanent"
)

type Customer struct {
	ID                      string     `gorm:"type:char(36);primaryKey" json:"id"`
	Token                   string     `gorm:"type:varchar(191);unique;not null" json:"token"`
	Phone                   *string    `gorm:"type:text" json:"phone"`
	PhoneVerifiedAt         *time.Time `json:"phone_verified_at"`
	EncryptedPII            *string    `gorm:"type:text" json:"encrypted_pii"`
	Status                  Status     `gorm:"type:enum('pending','verified','rejected');default:'pending'" json:"status"`
	KYCLevelAchieved        string     `gorm:"type:enum('none','tier_1','tier_2','tier_3');default:'none'" json:"kyc_level_achieved"`
	VerifiedAt              *time.Time `json:"verified_at"`
	IsBlacklisted           bool       `gorm:"default:false" json:"is_blacklisted"`
	Email                   string     `gorm:"type:text;not null" json:"email"`
	EmailHash               string     `gorm:"type:varchar(191);unique;not null" json:"email_hash"`
	DOB                     *time.Time `json:"dob"`
	Address                 *string    `gorm:"type:varchar(191)" json:"address"`
	AccessType              AccessType `gorm:"type:enum('temporary','permanent');default:'permanent'" json:"access_type"`
	FacialRecognitionPassed bool       `gorm:"default:false" json:"facial_recognition_passed"`
	CreatedAt               time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	// Permissions []Permission `gorm:"foreignKey:CustomerID" json:"permissions"`
	// Documents   []Document   `gorm:"foreignKey:CustomerID" json:"documents"`
	Identities []Identity `gorm:"foreignKey:CustomerID" json:"identities"`
}

func (Customer) TableName() string {
	return "customers"
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}
	return nil
}

func (customer *Customer) BeforeUpdate(tx *gorm.DB) (err error) {
	if customer.Status == StatusVerified && customer.VerifiedAt == nil {
		now := time.Now()
		customer.VerifiedAt = &now
	}
	if customer.Status != StatusVerified {
		customer.VerifiedAt = nil
	}
	return nil
}

func (customer *Customer) AfterFind(tx *gorm.DB) (err error) {
	if customer.Phone != nil {
		dec := Decrypt(*customer.Phone)
		customer.Phone = &dec
	}
	customer.Email = Decrypt(customer.Email)
	return nil
}

func (c *Customer) GetDecryptedPhone() *string {
	if c.Phone != nil {
		dec := Decrypt(*c.Phone)
		return &dec
	}
	return nil
}

func (c *Customer) GetDecryptedEmail() string {
	return Decrypt(c.Email)
}
