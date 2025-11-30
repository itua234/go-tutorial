package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrInvalidEncryptedFormat = errors.New("invalid encrypted format")
	ErrInvalidEncryptionKey   = errors.New("invalid encryption key")
)

// GenerateHexToken generates a cryptographically secure, 32-character hex token (16 bytes).
func GenerateHexToken() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// HashSHA256 computes the SHA256 hash of the value after normalizing it.
func HashSHA256(value string) string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// GenerateAPIKey creates a cryptographically secure, base64-encoded key with a prefix.
func GenerateAPIKey(prefix string) (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to read random bytes for key: %w", err)
	}
	return fmt.Sprintf("%s_%s", prefix, base64.URLEncoding.EncodeToString(b)), nil
}

// getEncryptionKey loads and validates the 32-byte (256-bit) AES key from the environment.
func getEncryptionKey() ([]byte, error) {
	keyHex := os.Getenv("ENCRYPTION_KEY")
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		return nil, ErrInvalidEncryptionKey
	}
	return key, nil
}

// Encrypt encrypts plaintext using AES-256 GCM and returns a hex-encoded IV and ciphertext separated by a colon.
func Encrypt(plaintext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create AES cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("could not create GCM: %w", err)
	}

	// Never reuse IVs/Nonces with GCM. Read a new one every time.
	iv := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not generate nonce: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), nil)

	// Combine IV and ciphertext for storage
	ivHex := hex.EncodeToString(iv)
	cipherHex := hex.EncodeToString(ciphertext)

	return fmt.Sprintf("%s:%s", ivHex, cipherHex), nil
}

// Decrypt decrypts a ciphertext produced by Encrypt.
func Decrypt(encrypted string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(encrypted, ":", 2)
	if len(parts) != 2 {
		return "", ErrInvalidEncryptedFormat
	}
	iv, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid IV hex: %w", err)
	}
	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext hex: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create AES cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("could not create GCM: %w", err)
	}

	// Validate Nonce/IV size
	if len(iv) != aesgcm.NonceSize() {
		return "", ErrInvalidEncryptedFormat
	}

	plaintext, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		// This error often indicates a failed authentication tag (tampering)
		return "", errors.New("decryption failed or data tampered")
	}

	return string(plaintext), nil
}
