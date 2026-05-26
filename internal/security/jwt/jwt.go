package jwtutil

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims holds JWT payload with a user identifier.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Config holds JWT signing parameters.
type Config struct {
	Secret        string
	TTL           time.Duration
	SecureCookies bool
}

// ConfigFromEnv builds a Config from environment variables.
// Requires JWT_SECRET. Optionally reads JWT_TTL (default 24h)
// and JWT_COOKIE_SECURE (default false).
func ConfigFromEnv() (*Config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	ttl := 24 * time.Hour
	if raw := os.Getenv("JWT_TTL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil {
			ttl = parsed
		}
	}

	return &Config{
		Secret:        secret,
		TTL:           ttl,
		SecureCookies: os.Getenv("JWT_COOKIE_SECURE") == "true",
	}, nil
}

// IssueToken creates a signed JWT for the given user.
func (c *Config) IssueToken(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("userID must not be nil")
	}

	now := time.Now()
	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(c.TTL)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(c.Secret))
}

// ValidateToken parses and validates a JWT, returning its claims.
func (c *Config) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token must not be empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, c.keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}

// SetAuthCookie issues a JWT and attaches it as an HttpOnly cookie.
func (c *Config) SetAuthCookie(w http.ResponseWriter, userID uuid.UUID) error {
	token, err := c.IssueToken(userID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(c.TTL.Seconds()),
		HttpOnly: true,
		Secure:   c.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// keyFunc returns the signing key after validating the algorithm.
func (c *Config) keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte(c.Secret), nil
}
