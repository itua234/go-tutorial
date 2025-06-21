package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestStatus string

const (
	RequestStatusInitiated     RequestStatus = "initiated"
	RequestStatusOTPPending    RequestStatus = "otp_pending"
	RequestStatusKYCProcessing RequestStatus = "kyc_processing"
	RequestStatusCompleted     RequestStatus = "completed"
	RequestStatusFailed        RequestStatus = "failed"
)

type KYCLevel string

const (
	KYCLevelTier1 KYCLevel = "tier_1"
	KYCLevelTier2 KYCLevel = "tier_2"
	KYCLevelTier3 KYCLevel = "tier_3"
)

type Request struct {
	ID                    string        `gorm:"type:char(36);primaryKey;unique" json:"id"`
	Reference             string        `gorm:"type:varchar(191);not null;unique" json:"reference"`
	RedirectURL           string        `gorm:"type:varchar(191);not null" json:"redirect_url"`
	KYCLevel              string        `gorm:"type:enum('tier_1','tier_2','tier_3');default:'tier_1'" json:"kyc_level"`
	BankAccountsRequested bool          `gorm:"default:false" json:"bank_accounts_requested"`
	EncryptedData         *string       `gorm:"type:text" json:"encrypted_data"`
	AllowURL              *string       `gorm:"type:varchar(191)" json:"allow_url"`
	KYCToken              string        `gorm:"type:varchar(191)" json:"kyc_token"`
	TokenExpiresAt        time.Time     `gorm:"not null" json:"token_expires_at"`
	CompanyID             string        `gorm:"type:char(36);not null" json:"company_id"`
	Company               *Company      `gorm:"foreignKey:CompanyID" json:"company"`
	Status                RequestStatus `gorm:"type:enum('initiated','otp_pending','kyc_processing','completed','failed');default:'initiated'" json:"status"`
	CreatedAt             time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Request) TableName() string {
	return "requests"
}

func (request *Request) BeforeCreate(tx *gorm.DB) (err error) {
	if request.ID == "" {
		request.ID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	if request.AllowURL == nil || *request.AllowURL == "" {
		generatedURL := "http://127.0.0.1:5173/" + request.KYCToken
		request.AllowURL = &generatedURL
	}

	return nil
}
