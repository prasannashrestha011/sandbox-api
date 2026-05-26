package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"main/internal/repository/model"
)

// RefreshTokenRepository defines persistence methods for refresh tokens.
//
// Typical usage:
// - Store: Create
//
// - Validate: FindValidByHash
//
// - Logout/rotation: RevokeByID / RevokeAllByUserID
//
// - Cleanup: DeleteExpired
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindValidByHash(ctx context.Context, tokenHash string, now time.Time) (*model.RefreshToken, error)
	RevokeByID(ctx context.Context, id uint) error
	RevokeAllByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context, now time.Time) (int64, error)
}

type gormRefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository returns a GORM-backed RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &gormRefreshTokenRepository{db: db}
}

func (r *gormRefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	if token == nil {
		return gorm.ErrInvalidData
	}
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).Create(token).Error
}

func (r *gormRefreshTokenRepository) FindValidByHash(ctx context.Context, tokenHash string, now time.Time) (*model.RefreshToken, error) {
	if tokenHash == "" {
		return nil, gorm.ErrInvalidData
	}

	var token model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = false AND expires_at > ?", tokenHash, now).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *gormRefreshTokenRepository) RevokeByID(ctx context.Context, id uint) error {
	if id == 0 {
		return gorm.ErrInvalidData
	}
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *gormRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	if userID == "" {
		return gorm.ErrInvalidData
	}
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func (r *gormRefreshTokenRepository) DeleteExpired(ctx context.Context, now time.Time) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("expires_at <= ?", now).
		Delete(&model.RefreshToken{})
	return res.RowsAffected, res.Error
}
