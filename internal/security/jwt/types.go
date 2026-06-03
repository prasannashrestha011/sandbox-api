package jwtutil

import (
	"main/internal/enums"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds JWT payload with a user identifier.
type Claims struct {
	UserID   string         `json:"user_id"`
	Role     string         `json:"role"`
	UserType enums.UserType `json:"user_type"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Config holds JWT signing parameters.
type Config struct {
	Secret        string
	TTL           time.Duration
	RefreshTTL    time.Duration
	SecureCookies bool
}
