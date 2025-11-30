package models

import (
	"confam-api/internal/crypto"
	"os"
	"testing"
)

func TestRequestTableName(t *testing.T) {
	request := Request{}
	expected := "requests"
	got := request.TableName()

	if got != expected {
		t.Errorf("TableName() = %v, want %v", got, expected)
	}
}

func TestBeforeCreate_GeneratesID(t *testing.T) {
	request := Request{}
	err := request.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("BeforeCreate() returned error: %v", err)
	}

	if request.ID == "" {
		t.Error("BeforeCreate() did not generate ID")
	}

	if len(request.ID) != 32 {
		t.Errorf("Generated ID length = %d, want 32", len(request.ID))
	}
}

func TestBeforeCreate_PreservesExistingID(t *testing.T) {
	existingID := "my-custom-id-123"
	request := Request{ID: existingID}
	err := request.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("BeforeCreate() returned error: %v", err)
	}

	if request.ID != existingID {
		t.Errorf("BeforeCreate() changed ID: got %v, want %v", request.ID, existingID)
	}
}

func TestBeforeCreate_GeneratesAllowURL_WhenNil(t *testing.T) {
	request := Request{
		KYCToken: "test-token-123",
	}
	err := request.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("BeforeCreate() returned error: %v", err)
	}

	if *request.AllowURL == "" {
		t.Error("BeforeCreate() did not generate AllowURL")
	}

	expectedURL := "http://127.0.0.1:5173/" + request.KYCToken
	if *request.AllowURL != expectedURL {
		t.Errorf("AllowURL = %v, want %v", *request.AllowURL, expectedURL)
	}
}

func TestBeforeCreate_PreservesExistingAllowURL(t *testing.T) {
	existingURL := "https://custom-url.com"
	request := Request{
		KYCToken: "test-token",
		AllowURL: &existingURL,
	}
	err := request.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("BeforeCreate() returned error: %v", err)
	}

	if *request.AllowURL != existingURL {
		t.Errorf("BeforeCreate() changed AllowURL: got %v, want %v", *request.AllowURL, existingURL)
	}
}

func TestBeforeCreate_GeneratesAllowURL_WhenEmptyString(t *testing.T) {
	emptyURL := ""
	request := Request{
		KYCToken: "test-token",
		AllowURL: &emptyURL,
	}
	err := request.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("BeforeCreate() returned error: %v", err)
	}

	expectedURL := "http://127.0.0.1:5173/" + request.KYCToken
	if *request.AllowURL != expectedURL {
		t.Errorf("AllowURL = %v, want %v", *request.AllowURL, expectedURL)
	}
}

func TestAfterFind_DecryptsData(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "7c633361cb709e1cf6ef0d68c914b5a5b1b540034331ab64702d2fd980dc7585")
	defer os.Unsetenv("ENCRYPTION_KEY")

	plainText := "sensitive-data"
	encrypted, _ := crypto.Encrypt(plainText)

	request := Request{EncryptedData: &encrypted}
	err := request.AfterFind(nil)
	if err != nil {
		t.Fatalf("AfterFind() returned error: %v", err)
	}

	if request.EncryptedData == nil {
		t.Fatal("AfterFind() set EncryptedData to nil")
	}

	if *request.EncryptedData != plainText {
		t.Errorf("Decrypted data = %v, want %v", *request.EncryptedData, plainText)
	}
}

func TestAfterFind_HandlesNilEncryptedData(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "7c633361cb709e1cf6ef0d68c914b5a5b1b540034331ab64702d2fd980dc7585")
	defer os.Unsetenv("ENCRYPTION_KEY")

	request := Request{EncryptedData: nil}
	err := request.AfterFind(nil)
	if err != nil {
		t.Fatalf("AfterFind() returned error: %v", err)
	}

	if request.EncryptedData != nil {
		t.Error("AfterFind() modified nil EncryptedData")
	}
}

func TestAfterFind_HandlesEmptyEncryptedData(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "7c633361cb709e1cf6ef0d68c914b5a5b1b540034331ab64702d2fd980dc7585")
	defer os.Unsetenv("ENCRYPTION_KEY")

	emptyString := ""
	request := Request{EncryptedData: &emptyString}
	err := request.AfterFind(nil)
	if err != nil {
		t.Fatalf("AfterFind() returned error: %v", err)
	}

	if *request.EncryptedData != emptyString {
		t.Errorf("AfterFind() modified empty EncryptedData")
	}
}
