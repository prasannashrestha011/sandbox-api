package dto

import (
	"time"

	"github.com/google/uuid"
)

// User is returned to callers and never contains a password.
type User struct {
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Fullname  string    `json:"fullname,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateUserInput captures registration input.
type CreateUserInput struct {
	Fullname string `json:"fullname,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"` // plain text, hashed in service before storage
}

// UpdateUserInput captures profile updates.
type UpdateUserInput struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Fullname string    `json:"fullname,omitempty"`
	Username string    `json:"username,omitempty"`
}

// UserCreate is used internally when storing a new user.
type UserCreate struct {
	Fullname     string `json:"fullname,omitempty"`
	Username     string `json:"username,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
}

// UserPasswordUpdate is used internally when updating the password.
type UserPasswordUpdate struct {
	UserID       uuid.UUID `json:"user_id,omitempty"`
	PasswordHash string    `json:"password_hash,omitempty"`
}

// LoginInput is used for authentication input.
type LoginInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// AuthResult is the authentication response payload.
type AuthResult struct {
	User  User   `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}
