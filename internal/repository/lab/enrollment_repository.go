package repository

import (
	"context"
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
	labGorm := mapper.EnrollmentToGorm(req)

	err := e.db.WithContext(ctx).Model(&labGorm).Create(&labGorm).Error
	if err != nil {
		return err
	}
	return nil

}

func (e *enrollmentRepository) GetEnrollment(ctx context.Context, userID string, labID string) (*models.LabEnrollment, error) {
	var enrollmentGorm lab_model.LabEnrollment
	err := e.db.WithContext(ctx).Where("user_id = ? AND lab_id = ?", userID, labID).First(&enrollmentGorm).Error
	if err != nil {
		return nil, err
	}
	return mapper.EnrollmentFromGorm(&enrollmentGorm), nil
}

func (e *enrollmentRepository) GetUserEnrollments(ctx context.Context, userID string) ([]models.LabEnrollment, error) {
	var enrollmentGorms []lab_model.LabEnrollment
	err := e.db.WithContext(ctx).Where("user_id = ?", userID).Find(&enrollmentGorms).Error
	if err != nil {
		return nil, err
	}
	return mapper.EnrollmentsFromGorm(&enrollmentGorms), nil
}

func (e *enrollmentRepository) UpdateEnrollmentProgress(ctx context.Context, req *models.LabEnrollment) error {
	labGorm := mapper.EnrollmentToGorm(req)

	err := e.db.WithContext(ctx).Model(&labGorm).Where("user_id = ? AND lab_id = ?", req.UserID, req.LabID).Updates(&labGorm).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *enrollmentRepository) DeleteEnrollment(ctx context.Context, userID string, labID string) error {
	var labGorm lab_model.LabEnrollment
	err := e.db.WithContext(ctx).Where("user_id = ? AND lab_id = ?", userID, labID).Delete(&labGorm).Error
	if err != nil {
		return err
	}
	return nil
}
