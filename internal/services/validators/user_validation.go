package services_validators

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"main/internal/repository/model"
)

func ValidateCreateUser(user *model.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}
	password := strings.TrimSpace(user.Password)
	if password == "" {
		return "", errors.New("password is required")
	}
	if len(password) > 72 {
		return "", errors.New("password must be 72 bytes or less")
	}
	if len([]rune(password)) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}

	return password, nil
}

func ValidateUpdateDetails(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.UserID == uuid.Nil {
		return errors.New("user id is required")
	}
	if strings.TrimSpace(user.Fullname) == "" && strings.TrimSpace(user.Username) == "" {
		return errors.New("no fields to update")
	}

	return nil
}

func ValidatePasswordUpdateInputs(id uuid.UUID, oldPassword, newPassword string) (string, string, error) {
	if id == uuid.Nil {
		return "", "", errors.New("user id is required")
	}
	oldPassword = strings.TrimSpace(oldPassword)
	if oldPassword == "" {
		return "", "", errors.New("old password is required")
	}
	newPassword = strings.TrimSpace(newPassword)
	if newPassword == "" {
		return "", "", errors.New("password is required")
	}
	if oldPassword == newPassword {
		return "", "", errors.New("new password must be different from old password")
	}
	if len(oldPassword) > 72 || len(newPassword) > 72 {
		return "", "", errors.New("password must be 72 bytes or less")
	}
	if len([]rune(newPassword)) < 8 {
		return "", "", errors.New("password must be at least 8 characters")
	}

	return oldPassword, newPassword, nil
}

func ValidatePasswordLookup(id uuid.UUID, password string) (string, error) {
	if id == uuid.Nil {
		return "", errors.New("user id is required")
	}
	password = strings.TrimSpace(password)
	if password == "" {
		return "", errors.New("password is required")
	}
	if len(password) > 72 {
		return "", errors.New("password must be 72 bytes or less")
	}

	return password, nil
}

func ValidateDeleteID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("user id is required")
	}
	return nil
}

func ValidateUsernameExistsInput(username string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return "", errors.New("username is required")
	}
	return username, nil
}
