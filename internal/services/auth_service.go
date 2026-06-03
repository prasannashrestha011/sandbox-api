package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"main/internal/domain"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/security/hashing"
	jwtutil "main/internal/security/jwt"
)

var ErrUnauthorized = errors.New("unauthorized")

// AuthService exposes authentication-related operations.
type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (*model.User, string, string, error)
	RefreshAccessToken(ctx context.Context, id uuid.UUID, refreshToken string) (string, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type authService struct {
	userrepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(userrepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &authService{
		userrepo:         userrepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *authService) Authenticate(ctx context.Context, username, password string) (*model.User, string, string, error) {

	username, password = strings.TrimSpace(username), strings.TrimSpace(password)

	user, err := s.userrepo.FindByUsername(ctx, username)
	if err != nil || user.Password == "" {
		return nil, "", "", domain.NotFoundError(err)
	}
	if err := hashing.ComparePasswordHash(user.Password, password); err != nil {
		return nil, "", "", domain.InvalidRequestError("invalid credentials", nil)
	}

	accessToken, refreshToken, err := jwtutil.JwtUtil.IssueToken(user.UserID, string(user.Role), user.UserType)
	if err != nil {
		return nil, "", "", domain.InternalError(errors.New("failed to issue token"))
	}

	if err := s.refreshTokenRepo.Create(ctx, &model.RefreshToken{
		UserID:    user.UserID.String(),
		TokenHash: hashing.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}); err != nil {
		return nil, "", "", domain.InternalError(errors.New("failed to store refresh token"))
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, id uuid.UUID, refreshToken string) (string, error) {

	stored, err := s.refreshTokenRepo.FindValidByHash(ctx, hashing.HashToken(refreshToken), time.Now())
	if err != nil || stored == nil || stored.UserID != id.String() {
		return "", postgres_error.MapError(err, "get refresh token by hash", "refresh_token")
	}

	user, err := s.userrepo.FindByID(ctx, id)
	if err != nil {
		return "", postgres_error.MapError(err, "get user by ID", "user")
	}

	accessToken, err := jwtutil.JwtUtil.IssueAccessToken(user.UserID, string(user.Role), user.UserType)
	if err != nil {
		return "", postgres_error.MapError(err, "issue access token", "access_token")
	}

	return accessToken, nil
}

func (s *authService) UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {

	valid, err := s.validatePassword(ctx, id, oldPassword)
	if err != nil {
		return err
	}
	if !valid {
		return &domain.AppError{Code: "INVALID_CREDENTIAL", Message: "old password is incorrect"}
	}

	hashed, _ := hashing.HashPassword(newPassword)
	return s.userrepo.UpdatePassword(ctx, id, hashed)
}

func (s *authService) UsernameExists(ctx context.Context, username string) (bool, error) {
	return s.userrepo.UsernameExists(ctx, username)
}

func (s *authService) validatePassword(ctx context.Context, id uuid.UUID, password string) (bool, error) {
	storedHash, err := s.userrepo.GetPasswordHash(ctx, id)
	if err != nil {
		return false, postgres_error.MapError(err, "get password hash by user ID", "user")
	}
	if err := hashing.ComparePasswordHash(storedHash, password); err != nil {
		return false, nil
	}

	return true, nil
}
