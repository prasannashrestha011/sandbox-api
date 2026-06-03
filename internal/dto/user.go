package dto

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Fullname  string    `json:"fullname,omitempty"`
	Username  string    `json:"username,omitempty"`
	Role      string    `json:"role,omitempty"` // "admin" or "user"
	UserType  string    `json:"type,omitempty"` // "internal" or "external"
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ### CreateUserInput captures registration input. ###
type CreateUserInput struct {
	Fullname string `json:"fullname,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"` // plain text, hashed in service before storage
	Role     string `json:"role,omitempty"`     // "admin" or "user"
	UserType string `json:"type,omitempty"`     // "internal" or "external"
}

func (r *CreateUserInput) Sanitize() {
	r.Fullname = strings.TrimSpace(r.Fullname)
	r.Username = strings.TrimSpace(r.Username)
	r.Role = strings.TrimSpace(r.Role)
}

func (r *CreateUserInput) Validate() error {
	r.Sanitize()
	var v ValidationErrors
	if strings.TrimSpace(r.Fullname) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "fullname",
			Message: "fullname is required",
		})
	}
	if strings.TrimSpace(r.Username) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "username",
			Message: "username is required",
		})
	}
	if strings.TrimSpace(r.Password) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "password",
			Message: "password is required",
		})
	} else if len(strings.TrimSpace(r.Password)) < 6 {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "password",
			Message: "password length should be at least six characters",
		})
	}
	if strings.TrimSpace(r.Role) != "admin" && strings.TrimSpace(r.Role) != "user" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "role",
			Message: "role must be admin or user",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

//###

// ### UpdateUserInput captures profile updates. ###
type UpdateUserInput struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Fullname string    `json:"fullname,omitempty"`
	Username string    `json:"username,omitempty"`
	Role     string    `json:"role,omitempty"` // "admin" or "user"
}

func (r *UpdateUserInput) Sanitize() {
	r.Fullname = strings.TrimSpace(r.Fullname)
	r.Username = strings.TrimSpace(r.Username)
	r.Role = strings.TrimSpace(r.Role)
}

func (r *UpdateUserInput) Validate() error {
	r.Sanitize()
	var v ValidationErrors
	if strings.TrimSpace(r.Fullname) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "fullname",
			Message: "fullname is required",
		})
	}
	if strings.TrimSpace(r.Username) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "username",
			Message: "username is required",
		})
	}
	if strings.TrimSpace(r.Role) != "admin" && strings.TrimSpace(r.Role) != "user" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "role",
			Message: "role must be admin or user",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

//###

// ### UserCreate is used internally when storing a new user. ###
type UserCreate struct {
	Fullname     string `json:"fullname,omitempty"`
	Username     string `json:"username,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	Role         string `json:"role,omitempty"` // "admin" or "user"
	UserType     string `json:"type,omitempty"` // "internal" or "external"
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

func (r *LoginInput) Sanitize() {
	r.Username = strings.TrimSpace(r.Username)
}

func (r *LoginInput) Validate() error {
	r.Sanitize()
	var v ValidationErrors
	if strings.TrimSpace(r.Username) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "username",
			Message: "username is required",
		})
	}
	if strings.TrimSpace(r.Password) == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "password",
			Message: "password is required",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

//###

// AuthResult is the authentication response payload.
type AuthResult struct {
	User  User   `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}
