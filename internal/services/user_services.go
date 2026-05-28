package services

import (
	"context"

	"github.com/google/uuid"

	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/security/hashing"
	services_validators "main/internal/services/validators"
)

// UserService exposes business operations for users (CRUD and profile updates).
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	UpdateDetails(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userrepo repository.UserRepository
}

// NewUserService returns a service backed by a UserRepository.
func NewUserService(userrepo repository.UserRepository) UserService {
	return &userService{userrepo: userrepo}
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

func (s *userService) List(ctx context.Context) ([]model.User, error) {
	return s.userrepo.List(ctx)
}

func (s *userService) UpdateDetails(ctx context.Context, user *model.User) error {
	if err := services_validators.ValidateUpdateDetails(user); err != nil {
		return err
	}
	return s.userrepo.UpdateDetails(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := services_validators.ValidateDeleteID(id); err != nil {
		return err
	}
	return s.userrepo.Delete(ctx, id)
}
