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

type SandboxInstanceRepository interface {
	Create(ctx context.Context, req *models.SandboxInstance) (*models.SandboxInstance, error)
	FindActiveSessionByUser(ctx context.Context, userID string, templateID string) (*models.SandboxInstance, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
}

func NewSandboxInstanceRepository(db *gorm.DB) SandboxInstanceRepository {
	return &sandboxRepository{db: db}
}

func (s *sandboxRepository) Create(ctx context.Context, req *models.SandboxInstance) (*models.SandboxInstance, error) {
	newinstance := mapper.InstanceToGorm(req)
	if err := s.db.WithContext(ctx).Model(&gormodel.SandboxInstance{}).Create(newinstance).Error; err != nil {
		return nil, err
	}
	return mapper.InstanceFromGorm(newinstance), nil
}

func (s *sandboxRepository) Delete(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&gormodel.SandboxInstance{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *sandboxRepository) FindActiveSessionByUser(
	ctx context.Context,
	userID string,
	templateID string,
) (*models.SandboxInstance, error) {

	var session gormodel.SandboxInstance

	err := s.db.WithContext(ctx).
		Where(
			"user_id = ? AND id = ? AND status = ? AND expires_at > NOW()",
			userID,
			templateID,
			"active",
		).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return mapper.InstanceFromGorm(&session), nil
}

func (s *sandboxRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := s.db.WithContext(ctx).Model(&gormodel.SandboxInstance{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
