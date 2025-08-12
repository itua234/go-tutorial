package utils

import "golang.org/x/crypto/bcrypt"

// ComparePasswordAndHash compares a plaintext password with a hashed password.
func ComparePasswordAndHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
