package repository

import (
	"context"

	"gorm.io/gorm"

	"main/internal/repository/mapper"
	lab_model "main/internal/repository/model/lab"
	"main/internal/services/models"
)

// LabRepository defines persistence methods for labs.
type LabRepository interface {
	Create(ctx context.Context, l *models.Lab) error
	FindByID(ctx context.Context, id string) (*models.Lab, error)
	Update(ctx context.Context, l *models.Lab) error
	Delete(ctx context.Context, id string) error
	// Note: You can add more methods here like FindAll, ListByTeacherID, etc.
}

type labRepository struct {
	db *gorm.DB
}

// NewLabRepository returns a GORM-backed LabRepository.
func NewLabRepository(db *gorm.DB) LabRepository {
	return &labRepository{db: db}
}

func (r *labRepository) Create(ctx context.Context, l *models.Lab) error {
	// Map service model to GORM model
	gormLab := mapper.LabToGorm(l)
	
	if err := r.db.WithContext(ctx).Create(gormLab).Error; err != nil {
		return err
	}
	
	return nil
}

func (r *labRepository) FindByID(ctx context.Context, id string) (*models.Lab, error) {
	var gormLab lab_model.Lab
	
	err := r.db.WithContext(ctx).
		Preload("Exercises").
		Preload("Tags").
		Preload("CreatedBy").
		First(&gormLab, "id = ?", id).Error
		
	if err != nil {
		return nil, err
	}
	
	// Map GORM model back to service model
	return mapper.LabFromGorm(&gormLab), nil
}

func (r *labRepository) Update(ctx context.Context, l *models.Lab) error {
	// Map service model to GORM model
	gormLab := mapper.LabToGorm(l)
	
	if err := r.db.WithContext(ctx).Save(gormLab).Error; err != nil {
		return err
	}
	
	return nil
}

func (r *labRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&lab_model.Lab{}, "id = ?", id).Error; err != nil {
		return err
	}
	
	return nil
}

