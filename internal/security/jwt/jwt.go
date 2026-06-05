package jwtutil

import (
	"errors"
	"main/internal/enums"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JwtUtil  *Config
	initOnce sync.Once
	initErr  error
)

// Init initializes the package-global JWT config exactly once.
// Subsequent calls are no-ops and will return the first init error (if any).
func Init(c *Config) error {
	initOnce.Do(func() {
		if c == nil {
			initErr = errors.New("jwt config is nil")
			return
		}
		if c.Secret == "" {
			initErr = errors.New("jwt secret must not be empty")
			return
		}
		if c.TTL <= 0 {
			c.TTL = 24 * time.Hour
		}
		if c.RefreshTTL <= 0 {
			c.RefreshTTL = 7 * c.TTL
		}
		JwtUtil = c
	})
	return initErr
}

// InitFromEnv loads config from env and initializes the package-global config once.
func InitFromEnv() error {
	c, err := ConfigFromEnv()
	if err != nil {
		return err
	}
	return Init(c)
}

func ConfigFromEnv() (*Config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT SECRET is not set")
	}
	ttl := 24 * time.Hour
	if raw := os.Getenv("JWT_TTL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil {
			ttl = parsed
		}
	}

	refreshTTL := ttl * 2
	if raw := os.Getenv("JWT_REFRESH_TTL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil {
			refreshTTL = parsed
		}
	}
	return &Config{
		Secret:        secret,
		TTL:           ttl,
		RefreshTTL:    refreshTTL,
		SecureCookies: os.Getenv("JWT_COOKIE_SECURE") == "true",
	}, nil
}

func (c *Config) IssueAccessToken(userID string, role string, userType enums.UserType) (string, error) {
	if userID == "" {
		return "", errors.New("userID must not be empty")
	}
	now := time.Now()
	return sign(Claims{
		UserID:   userID,
		Role:     role,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(c.TTL)),
		},
	}, c.Secret)
}

func (c *Config) IssueToken(userID string, role string, userType enums.UserType) (access, refresh string, err error) {
	if userID == "" {
		return "", "", errors.New("userID must not be empty")
	}
	now := time.Now()

	access, accessErr := sign(Claims{
		UserID:   userID,
		Role:     role,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(c.TTL)),
		},
	}, c.Secret)

	refresh, refreshErr := sign(RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(c.RefreshTTL)),
		},
	}, c.Secret)

	return access, refresh, errors.Join(accessErr, refreshErr)
}

func (c *Config) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token must not be empty")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc(c.Secret))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func (c *Config) SetAuthCookie(w http.ResponseWriter, accessToken, refreshToken string) {
	base := http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   c.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	}
	a := base
	a.Name, a.Value, a.MaxAge = "auth_token", accessToken, int(c.TTL.Seconds())

	r := base
	r.Name, r.Value, r.MaxAge = "refresh_token", refreshToken, int(c.RefreshTTL.Seconds())

	http.SetCookie(w, &a)
	http.SetCookie(w, &r)
}

// sign is the single internal helper for creating signed JWTs.
func sign(claims jwt.Claims, secret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
}

// keyFunc returns a jwt.Keyfunc closed over the provided secret.
func keyFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}
}
