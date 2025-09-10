package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NextOfKin struct {
	ID         string    `gorm:"type:char(36);primaryKey;unique" json:"id"`
	CustomerID string    `gorm:"type:char(36);not null" json:"customer_id"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"customer"`

	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	Email        string `gorm:"type:varchar(191)" json:"email"`
	Phone        string `gorm:"type:varchar(20);not null" json:"phone"`
	Relationship string `gorm:"type:varchar(50);not null" json:"relationship"`
	Address      string `gorm:"type:varchar(255)" json:"address"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (NextOfKin) TableName() string {
	return "next_of_kins"
}

func (n *NextOfKin) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == "" {
		n.ID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return nil
}
