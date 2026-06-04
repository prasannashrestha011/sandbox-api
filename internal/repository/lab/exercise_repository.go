package repository

import (
	"context"
	"main/internal/repository/mapper"
	lab_model "main/internal/repository/model/lab"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type ExerciseRepository interface {
	CreateExercise(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error)
	GetExerciseByID(ctx context.Context, id string) (*models.Exercise, error)
	UpdateExercise(ctx context.Context, exercise *models.Exercise) error
	DeleteExercise(ctx context.Context, id string) error
	ListExercisesByChapterID(ctx context.Context, chapterID string) ([]*models.Exercise, error)
}

type exerciseRepository struct {
	db *gorm.DB
}

// CreateExercise implements [ExerciseRepository].
func (e *exerciseRepository) CreateExercise(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	exerciseGorm := mapper.ExerciseToGorm(exercise)
	err := e.db.WithContext(ctx).Model(lab_model.Exercise{}).Create(exerciseGorm).Error
	if err != nil {
		return nil, err
	}
	exerciseModel := mapper.ExerciseFromGorm(exerciseGorm)
	return exerciseModel, nil
}

// DeleteExercise implements [ExerciseRepository].
func (e *exerciseRepository) DeleteExercise(ctx context.Context, id string) error {
	err := e.db.WithContext(ctx).Delete(&lab_model.Exercise{}, "id = ?", id).Error
	return err
}

// GetExerciseByID implements [ExerciseRepository].
func (e *exerciseRepository) GetExerciseByID(ctx context.Context, id string) (*models.Exercise, error) {
	var exerciseGorm lab_model.Exercise
	err := e.db.WithContext(ctx).Model(&lab_model.Exercise{}).Where("id = ?", id).First(&exerciseGorm).Error
	if err != nil {
		return nil, err
	}
	exerciseModel := mapper.ExerciseFromGorm(&exerciseGorm)
	return exerciseModel, nil
}

// ListExercisesByChapterID implements [ExerciseRepository].
func (e *exerciseRepository) ListExercisesByChapterID(ctx context.Context, chapterID string) ([]*models.Exercise, error) {
	var exerciseGorms []lab_model.Exercise
	err := e.db.WithContext(ctx).Model(&lab_model.Exercise{}).Where("chapter_id = ?", chapterID).Find(&exerciseGorms).Error
	if err != nil {
		return nil, err
	}
	exercises := make([]*models.Exercise, len(exerciseGorms))
	for i, exerciseGorm := range exerciseGorms {
		exercises[i] = mapper.ExerciseFromGorm(&exerciseGorm)
	}
	return exercises, nil
}

// UpdateExercise implements [ExerciseRepository].
func (e *exerciseRepository) UpdateExercise(ctx context.Context, exercise *models.Exercise) error {
	exerciseGorm := mapper.ExerciseToGorm(exercise)
	err := e.db.WithContext(ctx).Model(&lab_model.Exercise{}).Where("id = ?", exercise.ID).Updates(exerciseGorm).Error
	return err
}

func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}
