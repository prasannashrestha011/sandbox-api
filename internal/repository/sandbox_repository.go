package repository

import (
	"context"
	"main/internal/repository/mapper"
	gormodel "main/internal/repository/model"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type sandboxRepository struct {
	db *gorm.DB
}

type SandboxRepository interface {
	Create(ctx context.Context, req *models.SandboxSession) (*models.SandboxSession, error)
	FindActiveSessionByUserAndTemplate(ctx context.Context, userID string, templateID string) (*models.SandboxSession, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
}

func NewSandboxRepository(db *gorm.DB) SandboxRepository {
	return &sandboxRepository{db: db}
}

func (s *sandboxRepository) Create(ctx context.Context, req *models.SandboxSession) (*models.SandboxSession, error) {
	newsession := mapper.SessionToGorm(req)
	if err := s.db.WithContext(ctx).Model(&gormodel.SandboxSession{}).Create(newsession).Error; err != nil {
		return nil, err
	}
	return mapper.SessionFromGorm(newsession), nil
}

func (s *sandboxRepository) Delete(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&gormodel.SandboxSession{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *sandboxRepository) FindActiveSessionByUserAndTemplate(
	ctx context.Context,
	userID string,
	templateID string,
) (*models.SandboxSession, error) {

	var session gormodel.SandboxSession

	err := s.db.WithContext(ctx).
		Where(
			"user_id = ? AND template_id = ? AND status = ? AND expires_at > NOW()",
			userID,
			templateID,
			"active",
		).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return mapper.SessionFromGorm(&session), nil
}

func (s *sandboxRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := s.db.WithContext(ctx).Model(&gormodel.SandboxSession{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
