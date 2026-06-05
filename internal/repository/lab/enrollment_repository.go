package repository

import (
	"context"
	"errors"
	"main/internal/repository/mapper"
	lab_model "main/internal/repository/model/lab"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	EnrollUserToLab(ctx context.Context, req *models.LabEnrollment) error
	UpdateEnrollmentProgress(ctx context.Context, req *models.LabEnrollment) error
	GetEnrollment(ctx context.Context, userID, labID string) (*models.LabEnrollment, error)
	GetUserEnrollments(ctx context.Context, userID string) ([]models.LabEnrollment, error)
	DeleteEnrollment(ctx context.Context, userID, labID string) error
}
type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (e *enrollmentRepository) EnrollUserToLab(ctx context.Context, req *models.LabEnrollment) error {

	var existing lab_model.LabEnrollment
	err := e.db.WithContext(ctx).Unscoped().Where("enrolled_user_id=? AND lab_id=?", req.UserID, req.LabID).First(&existing).Error
	if err == nil {
		if !existing.DeletedAt.Valid {
			return gorm.ErrDuplicatedKey
		}
		return e.db.WithContext(ctx).Unscoped().Model(&existing).Update("deleted_at", nil).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	labGorm := mapper.EnrollmentToGorm(req)
	return e.db.WithContext(ctx).Create(&labGorm).Error
}

func (e *enrollmentRepository) GetEnrollment(ctx context.Context, userID string, labID string) (*models.LabEnrollment, error) {
	var enrollmentGorm lab_model.LabEnrollment
	err := e.db.WithContext(ctx).Where("enrolled_user_id = ? AND lab_id = ?", userID, labID).First(&enrollmentGorm).Error
	if err != nil {
		return nil, err
	}
	return mapper.EnrollmentFromGorm(&enrollmentGorm), nil
}

func (e *enrollmentRepository) GetUserEnrollments(ctx context.Context, userID string) ([]models.LabEnrollment, error) {
	var enrollmentGorms []lab_model.LabEnrollment
	err := e.db.WithContext(ctx).Where("enrolled_user_id = ?", userID).Find(&enrollmentGorms).Error
	if err != nil {
		return nil, err
	}
	return mapper.EnrollmentsFromGorm(&enrollmentGorms), nil
}

func (e *enrollmentRepository) UpdateEnrollmentProgress(ctx context.Context, req *models.LabEnrollment) error {
	labGorm := mapper.EnrollmentToGorm(req)

	err := e.db.WithContext(ctx).Model(&labGorm).Where("enrolled_user_id = ? AND lab_id = ?", req.UserID, req.LabID).Updates(&labGorm).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *enrollmentRepository) DeleteEnrollment(ctx context.Context, userID string, labID string) error {
	// Execute the delete operation
	result := e.db.WithContext(ctx).
		Where("enrolled_user_id = ? AND lab_id = ?", userID, labID).
		Delete(&lab_model.LabEnrollment{})

	// 1. Check for actual database errors (connection issues, syntax, etc.)
	if result.Error != nil {
		return result.Error
	}

	// 2. Check if any rows were actually soft-deleted
	// If 0 rows were affected, it means the record didn't exist OR was already soft-deleted.
	if result.RowsAffected == 0 {
		// Return GORM's standard not found error (or your own custom error)
		return gorm.ErrRecordNotFound
	}

	return nil
}
