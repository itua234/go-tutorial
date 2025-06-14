package models

import (
	"time"
)

type IdentityType string

const (
	IdentityTypeNIN IdentityType = "NIN"
	IdentityTypeBVN IdentityType = "BVN"
)

type IdentityStatus string

const (
	IdentityStatusPending   IdentityStatus = "pending"
	IdentityStatusVerified  IdentityStatus = "verified"
	IdentityStatusRejected  IdentityStatus = "rejected"
	IdentityStatusExpired   IdentityStatus = "expired"
	IdentityStatusRevoked   IdentityStatus = "revoked"
	IdentityStatusSuspended IdentityStatus = "suspended"
)

func decryptString(encryptedValue string) string {
	return "DECRYPTED_" + encryptedValue
}

type Identity struct {
	ID                   string         `gorm:"type:char(36);primaryKey;unique" json:"id"`
	CustomerID           string         `gorm:"type:char(36);not null" json:"customer_id"`
	Customer             *Customer      `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"customer"`
	Type                 IdentityType   `gorm:"type:enum('NIN', 'BVN');not null" json:"type"`
	Value                string         `gorm:"type:text;not null" json:"value_encrypted"`
	Status               IdentityStatus `gorm:"type:enum('pending', 'verified', 'rejected', 'expired', 'revoked', 'suspended');default:'pending'" json:"status"`
	Verified             bool           `gorm:"not null;default:false" json:"verified"`
	VerificationProvider *string        `gorm:"type:varchar(255)" json:"verification_provider"`
	ProviderReference    *string        `gorm:"type:varchar(255)" json:"provider_reference"`
	VerifiedAt           *time.Time     `gorm:"type:timestamp;column:verified_at" json:"verified_at"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Identity) TableName() string {
	return "identities"
}

func (identity *Identity) GetDecryptedValue() string {
	if identity.Value == "" {
		return ""
	}
	return decryptString(identity.Value)
}

func (identity *Identity) IsShareable() bool {
	return true
}
