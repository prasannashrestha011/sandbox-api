package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/security/hashing"
	jwtutil "main/internal/security/jwt"
	services_validators "main/internal/services/validators"
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
	fail := func(msg string) (*model.User, string, string, error) {
		return nil, "", "", errors.New(msg)
	}

	username, password = strings.TrimSpace(username), strings.TrimSpace(password)
	if username == "" || password == "" {
		return fail("username and password are required")
	}

	user, err := s.userrepo.FindByUsername(ctx, username)
	if err != nil || user.Password == "" {
		return fail("invalid credentials")
	}
	if err := hashing.ComparePasswordHash(user.Password, password); err != nil {
		return fail("invalid credentials")
	}

	accessToken, refreshToken, err := jwtutil.JwtUtil.IssueToken(user.UserID, string(user.Role))
	if err != nil {
		return fail("failed to issue token")
	}

	if err := s.refreshTokenRepo.Create(ctx, &model.RefreshToken{
		UserID:    user.UserID.String(),
		TokenHash: hashing.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}); err != nil {
		return fail("failed to store refresh token")
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, id uuid.UUID, refreshToken string) (string, error) {
	fail := func(err error) (string, error) { return "", err }

	refreshToken, err := services_validators.ValidateRefreshAccessTokenInputs(id, refreshToken)
	if err != nil {
		return fail(err)
	}

	stored, err := s.refreshTokenRepo.FindValidByHash(ctx, hashing.HashToken(refreshToken), time.Now())
	if err != nil || stored == nil || stored.UserID != id.String() {
		return fail(ErrUnauthorized)
	}

	user, err := s.userrepo.FindByID(ctx, id)
	if err != nil {
		return fail(ErrUnauthorized)
	}

	accessToken, err := jwtutil.JwtUtil.IssueAccessToken(user.UserID, string(user.Role))
	if err != nil {
		return fail(errors.New("failed to issue token"))
	}

	return accessToken, nil
}

func (s *authService) UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
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

func (s *authService) UsernameExists(ctx context.Context, username string) (bool, error) {
	username, err := services_validators.ValidateUsernameExistsInput(username)
	if err != nil {
		return false, err
	}
	return s.userrepo.UsernameExists(ctx, username)
}

func (s *authService) validatePassword(ctx context.Context, id uuid.UUID, password string) (bool, error) {
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
