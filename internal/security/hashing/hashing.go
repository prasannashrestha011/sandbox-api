package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt with the default cost.
func HashPassword(password string) (string, error) {
	return HashPasswordWithCost(password, bcrypt.DefaultCost)
}

func HashPasswordWithCost(password string, cost int) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}
	if cost <= 0 {
		return "", errors.New("invalid bcrypt cost")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswordHash(passwordHash, password string) error {
	if passwordHash == "" {
		return errors.New("password hash is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func VerifyPassword(passwordHash, password string) bool {
	return ComparePasswordHash(passwordHash, password) == nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
