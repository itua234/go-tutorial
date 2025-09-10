package seeders

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"confam-api/models"
)

func Seed(db *gorm.DB, redisClient *redis.Client) error {
	// 1. Hash the password
	password := "Reckless@3030"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. Create company
	company := models.Company{
		ID:                   uuid.New().String(),
		Name:                 "Tech Innovators Inc.",
		Logo:                 ptr("https://example.com/logo1.png"),
		Email:                "johndoe2@example.com",
		Password:             ptr(string(hashedPassword)),
		Domain:               ptr("techinnovators.com"),
		WebhookURL:           ptr("https://webhooks.techinnovators.com"),
		Verified:             true,
		NotificationsEnabled: true,
		EmailVerifiedAt:      ptrTime(time.Now()),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := db.Create(&company).Error; err != nil {
		return err
	}

	// 3. Create app
	testPub := "pk_test_" + uuid.New().String()
	livePub := "pk_live_" + uuid.New().String()
	app := models.App{
		ID:            uuid.New().String(),
		Name:          "Sample App",
		DisplayName:   "Sample Application",
		Logo:          ptr("https://example.com/logo.png"),
		TestPublicKey: testPub,
		LivePublicKey: livePub,
		Mode:          "SANDBOX",
		Status:        "ACTIVE",
		WebhookURL:    ptr("https://example.com/webhook"),
		CompanyID:     company.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&app).Error; err != nil {
		return err
	}

	// 4. Store secret keys in Redis
	secretSalt := os.Getenv("SECRET_SALT")
	if secretSalt == "" {
		secretSalt = "default_secret_salt"
	}
	testSecret := "sk_test_" + hmacHex(secretSalt, testPub)
	liveSecret := "sk_live_" + hmacHex(secretSalt, livePub)

	fmt.Println("Test Secret:", testSecret)
	fmt.Println("Live Secret:", liveSecret)
	fmt.Println("App ID:", app.ID)

	ctx := context.Background()
	if err := redisClient.Set(ctx, "secret:"+testSecret, app.ID, 0).Err(); err != nil {
		return err
	}
	if err := redisClient.Set(ctx, "secret:"+liveSecret, app.ID, 0).Err(); err != nil {
		return err
	}

	return nil
}

// Helper functions
func ptr(s string) *string           { return &s }
func ptrTime(t time.Time) *time.Time { return &t }
func hmacHex(secret, msg string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}
