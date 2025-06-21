package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"crypto/sha256"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"os"
)

// GenerateToken generates a 32-character hex token (16 bytes = 32 hex chars)
func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// handle error as needed, here we just panic for simplicity
		panic(err)
	}
	return hex.EncodeToString(b)
}

func HashFunction(value string) string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	hash := sha256.Sum256([]byte(normalized))
	hexStr := hex.EncodeToString(hash[:])
	return hexStr
}

func Encrypt(plaintext string) string {
	keyHex := os.Getenv("ENCRYPTION_KEY")
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		panic("invalid encryption key")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	iv := make([]byte, aesgcm.NonceSize())
	io.ReadFull(rand.Reader, iv)

	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(iv) + ":" + hex.EncodeToString(ciphertext)
}

// Decrypt decrypts ciphertext produced by Encrypt
func Decrypt(encrypted string) (string, error) {
	keyHex := os.Getenv("ENCRYPTION_KEY")
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		return "", errors.New("invalid encryption key")
	}

	parts := strings.SplitN(encrypted, ":", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid encrypted format")
	}
	iv, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func generateApikey(prefix string) string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%s_%s", prefix, base64.URLEncoding.EncodeToString(b))
}
