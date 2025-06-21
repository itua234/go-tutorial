package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKey struct {
	ID        string     `gorm:"type:char(36);primaryKey" json:"id"`
	Key       string     `gorm:"unique;not null"`
	Type      string     `json:"type"`
	AppID     string     `gorm:"type:char(36);not null" json:"app_id"`
	App       *App       `gorm:"foreignKey:AppID" json:"app"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	//DeletedAt gorm.DeletedAt `gorm:"index"` // Enables soft deletes
}

func (APIKey) TableName() string {
	return "api_keys"
}

func (api_key *APIKey) BeforeCreate(tx *gorm.DB) (err error) {
	if api_key.ID == "" {
		api_key.ID = uuid.New().String()
	}
	return nil
}
