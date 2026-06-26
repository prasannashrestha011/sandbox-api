package repository

import (
	"context"
	"main/internal/repository/mapper"
	gormodel "main/internal/repository/model"
	"main/internal/services/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sandboxRepository struct {
	db *gorm.DB
}

type SandboxInstanceRepository interface {
	Create(ctx context.Context, req *models.SandboxInstance) (*models.SandboxInstance, error)
	Acquire(ctx context.Context, lang string) (*models.SandboxInstance, error)
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

func (s *sandboxRepository) Acquire(ctx context.Context, lang string) (*models.SandboxInstance, error) {
	var instance gormodel.SandboxInstance
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("lang =? AND status= ?", lang, "idle").First(&instance).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	instance.Status = "busy"
	if err := tx.Save(&instance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return mapper.InstanceFromGorm(&instance), nil
}

func (s *sandboxRepository) Delete(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&gormodel.SandboxInstance{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *sandboxRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := s.db.WithContext(ctx).Model(&gormodel.SandboxInstance{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
