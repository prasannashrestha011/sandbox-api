package services

import (
	"context"

	"github.com/google/uuid"

	postgres_error "main/internal/infra/postgres"
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

	err := services_validators.ValidateUserTypeAndRole(user)
	if err != nil {
		return err
	}
	hashed, err := hashing.HashPassword(user.Password)
	if err != nil {
		return postgres_error.MapError(err, "create user", "user")
	}
	user.Password = hashed

	err = s.userrepo.Create(ctx, user)
	if err != nil {
		return postgres_error.MapError(err, "create user", "user").Err
	}
	return nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userrepo.FindByID(ctx, id)
	if err != nil {
		return nil, postgres_error.MapError(err, "get user by id", "user")
	}
	return user, nil
}

func (s *userService) List(ctx context.Context) ([]model.User, error) {
	users, err := s.userrepo.List(ctx)
	if err != nil {
		return nil, postgres_error.MapError(err, "list users", "user")
	}
	return users, nil
}

func (s *userService) UpdateDetails(ctx context.Context, user *model.User) error {
	err := s.userrepo.UpdateDetails(ctx, user)
	return postgres_error.MapError(err, "update user details", "user")

}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return postgres_error.MapError(s.userrepo.Delete(ctx, id), "delete user", "user")
}
