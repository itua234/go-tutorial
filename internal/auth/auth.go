package auth

import "golang.org/x/crypto/bcrypt"

func ComparePasswordAndHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
