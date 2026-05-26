package services

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/security/hashing"
	services_validators "main/internal/services/validators"
)

// UserService exposes business operations for users.
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Authenticate(ctx context.Context, username, password string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	UpdateDetails(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	Delete(ctx context.Context, id uuid.UUID) error
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type userService struct {
	userrepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

// NewUserService returns a service backed by a UserRepository.
func NewUserService(userrepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) UserService {
	return &userService{
		userrepo:         userrepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	password, err := services_validators.ValidateCreateUser(user)
	if err != nil {
		return err
	}

	hashed, err := hashing.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashed

	return s.userrepo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userrepo.FindByID(ctx, id)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userrepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	if err := hashing.ComparePasswordHash(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) List(ctx context.Context) ([]model.User, error) {
	return s.userrepo.List(ctx)
}

func (s *userService) UpdateDetails(ctx context.Context, user *model.User) error {
	if err := services_validators.ValidateUpdateDetails(user); err != nil {
		return err
	}

	return s.userrepo.UpdateDetails(ctx, user)
}

func (s *userService) UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
	oldPassword, newPassword, err := services_validators.ValidatePasswordUpdateInputs(id, oldPassword, newPassword)
	if err != nil {
		return err
	}

	valid, err := s.validatePassword(ctx, id, oldPassword)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid old password")
	}

	hashed, err := hashing.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.userrepo.UpdatePassword(ctx, id, hashed)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := services_validators.ValidateDeleteID(id); err != nil {
		return err
	}
	return s.userrepo.Delete(ctx, id)
}

func (s *userService) UsernameExists(ctx context.Context, username string) (bool, error) {
	username, err := services_validators.ValidateUsernameExistsInput(username)
	if err != nil {
		return false, err
	}
	return s.userrepo.UsernameExists(ctx, username)
}

func (s *userService) validatePassword(ctx context.Context, id uuid.UUID, password string) (bool, error) {
	password, err := services_validators.ValidatePasswordLookup(id, password)
	if err != nil {
		return false, err
	}

	storedHash, err := s.userrepo.GetPasswordHash(ctx, id)
	if err != nil {
		return false, err
	}
	if storedHash == "" {
		return false, errors.New("stored password is empty")
	}

	if err := hashing.ComparePasswordHash(storedHash, password); err != nil {
		return false, nil
	}

	return true, nil
}
