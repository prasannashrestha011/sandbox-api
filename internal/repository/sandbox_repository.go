package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"main/internal/enums"
	"main/internal/repository/model"
)

type dockerRepository struct {
	db *gorm.DB
}

// SandboxRepository defines persistence methods for sandbox sessions.
type SandboxRepository interface {
	Create(ctx context.Context, sandbox *model.Sandbox) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Sandbox, error)
	FindBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Sandbox, error)
	ListByUserID(ctx context.Context, userID string) ([]model.Sandbox, error)
	UpdateStatus(ctx context.Context, containerID string, status enums.SandboxState) error
	Delete(ctx context.Context, containerID string) error
}

// NewSandboxRepository returns a GORM-backed SandboxRepository.
func NewSandboxRepository(db *gorm.DB) SandboxRepository {
	return &dockerRepository{db: db}
}

func (r *dockerRepository) Create(ctx context.Context, sandbox *model.Sandbox) error {
	return r.db.WithContext(ctx).Model(&model.Sandbox{}).Omit("Image").Create(sandbox).Error
}

func (r *dockerRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Sandbox, error) {
	var sandbox model.Sandbox
	err := r.db.WithContext(ctx).First(&sandbox, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sandbox, nil
}

func (r *dockerRepository) FindBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Sandbox, error) {
	var sandbox model.Sandbox
	err := r.db.WithContext(ctx).First(&sandbox, "session_id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &sandbox, nil
}

func (r *dockerRepository) ListByUserID(ctx context.Context, userID string) ([]model.Sandbox, error) {
	var sandboxes []model.Sandbox
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&sandboxes).Error
	if err != nil {
		return nil, err
	}
	return sandboxes, nil
}

func (r *dockerRepository) UpdateStatus(ctx context.Context, containerID string, status enums.SandboxState) error {
	return r.db.WithContext(ctx).Model(&model.Sandbox{}).Where("container_id = ?", containerID).Update("status", status).Error
}

func (r *dockerRepository) Delete(ctx context.Context, containerID string) error {
	return r.db.WithContext(ctx).Delete(&model.Sandbox{}, "container_id = ?", containerID).Error
}
