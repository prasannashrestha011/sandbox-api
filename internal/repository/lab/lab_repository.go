package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"main/internal/repository/mapper"
	lab_model "main/internal/repository/model/lab"
	"main/internal/services/models"
)

type labRepository struct {
	db *gorm.DB
}

// LabRepository defines persistence methods for labs.
type LabRepository interface {
	Create(ctx context.Context, l *models.Lab) (*models.Lab, error)
	FindByID(ctx context.Context, id string) (*models.Lab, error)
	Update(ctx context.Context, l *models.Lab) error
	Delete(ctx context.Context, id string) error
	// Note: You can add more methods here like FindAll, ListByTeacherID, etc.
}

// NewLabRepository returns a GORM-backed LabRepository.
func NewLabRepository(db *gorm.DB) LabRepository {
	return &labRepository{db: db}
}

func (r *labRepository) Create(ctx context.Context, l *models.Lab) (*models.Lab, error) {
	gormLab := mapper.LabToGorm(l)

	if err := r.db.WithContext(ctx).Omit(clause.Associations).Save(gormLab).Error; err != nil {
		return nil, err
	}

	// 2. Explicitly manage the tags inside the junction table.
	// Ensure tag records exist first to avoid foreign key constraint failures,
	// then replace the association so GORM updates the "lab_tags" join table.
	if len(gormLab.Tags) > 0 {
		// Upsert each tag record (create if missing).
		for i := range gormLab.Tags {
			tag := &gormLab.Tags[i]
			// If an ID is present, try to find by ID first.
			if tag.ID != "" {
				var existing lab_model.Tag
				if err := r.db.WithContext(ctx).First(&existing, "id = ?", tag.ID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				}
			} else {
				// No ID: find or create by name to avoid duplicates.
				if err := r.db.WithContext(ctx).Where(lab_model.Tag{Name: tag.Name}).FirstOrCreate(tag).Error; err != nil {
					return nil, err
				}
			}
		}

		if err := r.db.WithContext(ctx).Model(gormLab).Association("Tags").Replace(gormLab.Tags); err != nil {
			return nil, err
		}
	}

	serviceLab := mapper.LabFromGorm(gormLab)
	return serviceLab, nil
}

func (r *labRepository) FindByID(ctx context.Context, id string) (*models.Lab, error) {
	var gormLab lab_model.Lab

	err := r.db.WithContext(ctx).
		Preload("Chapters").
		Preload("Chapters.Exercises").
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
