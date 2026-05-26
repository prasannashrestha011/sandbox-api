package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"main/internal/repository/model"
)

// UserRepository defines persistence methods for users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	UpdateDetails(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	Delete(ctx context.Context, id uuid.UUID) error
	UsernameExists(ctx context.Context, username string) (bool, error)
	GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return gorm.ErrInvalidData
	}

	return r.db.WithContext(ctx).Model(&model.User{}).Create(user).Error
}

func (r *gormUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var dbUser model.User
	err := r.db.WithContext(ctx).First(&dbUser, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dbUser, nil
}

func (r *gormUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var dbUser model.User
	err := r.db.WithContext(ctx).Where("LOWER(username) = LOWER(?)", username).First(&dbUser).Error
	if err != nil {
		return nil, err
	}
	return &dbUser, nil
}

func (r *gormUserRepository) List(ctx context.Context) ([]model.User, error) {
	var dbUsers []model.User
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&dbUsers).Error
	if err != nil {
		return nil, err
	}
	return dbUsers, nil
}

func (r *gormUserRepository) UpdateDetails(ctx context.Context, user *model.User) error {
	if user == nil {
		return gorm.ErrInvalidData
	}
	return r.db.WithContext(ctx).Model(&model.User{}).Where("user_id = ?", user.UserID).Updates(user).Error
}

func (r *gormUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if id == uuid.Nil || passwordHash == "" {
		return gorm.ErrInvalidData
	}
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("user_id = ?", id).
		Update("password", passwordHash).Error
}

func (r *gormUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "user_id = ?", id).Error
}

func (r *gormUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("LOWER(username) = LOWER(?)", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *gormUserRepository) GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	var dbUser model.User
	err := r.db.WithContext(ctx).Select("password").First(&dbUser, "user_id = ?", id).Error
	if err != nil {
		return "", err
	}
	return dbUser.Password, nil
}
