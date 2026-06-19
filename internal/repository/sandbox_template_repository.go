package repository

import (
	"context"

	"gorm.io/gorm"

	"main/internal/repository/mapper"
	gormodel "main/internal/repository/model"
	"main/internal/services/models"
)

type sandboxTemplateRepository struct {
	db *gorm.DB
}

// SandboxTemplateRepository defines persistence methods for sandbox templates.
type SandboxTemplateRepository interface {
	Create(ctx context.Context, req *models.SandboxTemplate) (*models.SandboxTemplate, error)
	FindByID(ctx context.Context, id string) (*models.SandboxTemplate, error)
	ListByUserID(ctx context.Context, userID string) ([]models.SandboxTemplate, error)
	UpdateDetails(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, containerID string) error
}

// NewSandboxTemplateRepository returns a GORM-backed SandboxTemplateRepository.
func NewSandboxTemplateRepository(db *gorm.DB) SandboxTemplateRepository {
	return &sandboxTemplateRepository{db: db}
}

func (r *sandboxTemplateRepository) Create(ctx context.Context, req *models.SandboxTemplate) (*models.SandboxTemplate, error) {
	sandbox := mapper.TemplateToGorm(req)
	err := r.db.WithContext(ctx).Model(&gormodel.SandboxTemplate{}).Omit("Image").Create(sandbox).Error
	if err != nil {
		return nil, err
	}
	return mapper.TemplateFromGorm(sandbox), nil
}

func (r *sandboxTemplateRepository) FindByID(ctx context.Context, id string) (*models.SandboxTemplate, error) {
	var sandbox gormodel.SandboxTemplate
	err := r.db.WithContext(ctx).Preload("Image").First(&sandbox, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return mapper.TemplateFromGorm(&sandbox), nil
}

func (r *sandboxTemplateRepository) ListByUserID(ctx context.Context, userID string) ([]models.SandboxTemplate, error) {
	var sandboxes []gormodel.SandboxTemplate
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&sandboxes).Error
	if err != nil {
		return nil, err
	}
	return mapper.TemplateListFromGorm(sandboxes), nil
}
func (r *sandboxTemplateRepository) UpdateDetails(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&gormodel.SandboxTemplate{}).
		Where("id = ?", id).
		Select("MemoryLimit", "PidsLimit", "CPULimit", "NetworkMode", "SessionTimeout", "ExecTimeout").
		Updates(updates).
		Error
}

func (r *sandboxTemplateRepository) Delete(ctx context.Context, containerID string) error {
	return r.db.WithContext(ctx).Delete(&gormodel.SandboxTemplate{}, "container_id = ?", containerID).Error
}
