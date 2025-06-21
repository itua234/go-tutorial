package models

import (
	client "confam-api/utils"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppMode string

const (
	AppModeSandbox AppMode = "SANDBOX"
	AppModeLive    AppMode = "LIVE"
)

var SecretSalt = func() string {
	if s := os.Getenv("SECRET_SALT"); s != "" {
		return s
	}
	return "default_secret_salt"
}()

type App struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	DisplayName   string    `gorm:"type:varchar(255);not null" json:"display_name"`
	Logo          *string   `gorm:"type:varchar(255)" json:"logo"`
	TestPublicKey string    `gorm:"type:varchar(64);unique;not null" json:"test_public_key"`
	LivePublicKey string    `gorm:"type:varchar(64);unique;not null" json:"live_public_key"`
	Mode          string    `gorm:"type:enum('SANDBOX','LIVE');default:SANDBOX" json:"mode"`
	Status        string    `gorm:"type:enum('ACTIVE','INACTIVE');default:ACTIVE" json:"status"`
	WebhookURL    *string   `gorm:"type:varchar(255)" json:"webhook_url"`
	CompanyID     string    `gorm:"type:char(36);not null" json:"company_id"`
	Company       *Company  `gorm:"foreignKey:CompanyID" json:"company"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (App) TableName() string {
	return "apps"
}

func (app *App) TestSecretKey() string {
	h := hmac.New(sha256.New, []byte(SecretSalt))
	h.Write([]byte(app.TestPublicKey))
	return "sk_test_" + hex.EncodeToString(h.Sum(nil))
}

func (app *App) LiveSecretKey() string {
	h := hmac.New(sha256.New, []byte(SecretSalt))
	h.Write([]byte(app.LivePublicKey))
	return "sk_live_" + hex.EncodeToString(h.Sum(nil))
}

func (app *App) BeforeCreate(tx *gorm.DB) (err error) {
	if app.ID == "" {
		app.ID = uuid.New().String()
	}
	if app.TestPublicKey == "" {
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			return errors.New("failed to generate random bytes for test public key")
		}
		app.TestPublicKey = "pk_test_" + hex.EncodeToString(bytes)
	}
	if app.LivePublicKey == "" {
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			return errors.New("failed to generate random bytes for live public key")
		}
		app.LivePublicKey = "pk_live_" + hex.EncodeToString(bytes)
	}
	return nil
}

// AfterCreate GORM hook to store secret keys in external client (pseudo-code)
func (app *App) AfterCreate(tx *gorm.DB) (err error) {
	testSecret := app.TestSecretKey()
	liveSecret := app.LiveSecretKey()
	fmt.Println("Test Secret:", testSecret)
	fmt.Println("Live Secret:", liveSecret)
	ctx := context.Background()
	err = client.RedisClient.Set(
		ctx,
		fmt.Sprintf("secret:%s", testSecret),
		app.ID,
		0, // 0 means no expiration
	).Err()
	if err != nil {
		return err
	}

	err = client.RedisClient.Set(
		ctx,
		fmt.Sprintf("secret:%s", liveSecret),
		app.ID,
		0,
	).Err()
	if err != nil {
		return err
	}

	return nil
}
