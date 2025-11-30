package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Webhook struct {
	ID            string `gorm:"type:char(36);primaryKey;unique" json:"id"`
	EventType     string
	Payload       string
	TargetURL     string
	Status        string `gorm:"type:enum('pending', 'success', 'failed');default:'pending'" json:"status"`
	Attempts      int
	LastAttemptAt *time.Time
	ErrorMessage  *string
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Webhook) TableName() string {
	return "webhooks"
}

func (webhook *Webhook) BeforeCreate(tx *gorm.DB) (err error) {
	if webhook.ID == "" {
		webhook.ID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return nil
}
