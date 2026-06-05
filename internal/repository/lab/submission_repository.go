package repository

import (
	"context"
	"main/internal/repository/mapper"
	lab_model "main/internal/repository/model/lab"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	CreateSubmission(ctx context.Context, req *models.Submission) (*models.Submission, error)
	GetSubmissionsByExercise(ctx context.Context, exerciseID string) ([]models.Submission, error)
	GetSubmissionsByUser(ctx context.Context, userID string) ([]models.Submission, error)
	GetSubmissionByID(ctx context.Context, submissionID string) (*models.Submission, error)
	UpdateSubmission(ctx context.Context, req *models.Submission) (*models.Submission, error)
	DeleteSubmission(ctx context.Context, submissionID string) error
}

type submissionRepository struct {
	db *gorm.DB // Assuming you have a GORM DB instance
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (s *submissionRepository) CreateSubmission(ctx context.Context, req *models.Submission) (*models.Submission, error) {
	submission := mapper.SubmissionToGorm(req)
	err := s.db.WithContext(ctx).Model(&lab_model.Submission{}).Create(submission).Error
	if err != nil {
		return nil, err
	}
	return mapper.SubmissionFromGorm(submission), nil
}

func (s *submissionRepository) DeleteSubmission(ctx context.Context, submissionID string) error {
	return s.db.WithContext(ctx).Delete(&lab_model.Submission{}, "id = ?", submissionID).Error
}

// GetSubmissionByID implements [SubmissionRepository].
func (s *submissionRepository) GetSubmissionByID(ctx context.Context, submissionID string) (*models.Submission, error) {
	var submission lab_model.Submission
	err := s.db.WithContext(ctx).Model(&lab_model.Submission{}).Where("id = ?", submissionID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return mapper.SubmissionFromGorm(&submission), nil
}

func (s *submissionRepository) GetSubmissionsByExercise(ctx context.Context, exerciseID string) ([]models.Submission, error) {
	var submissions []lab_model.Submission
	err := s.db.WithContext(ctx).Model(&lab_model.Submission{}).Where("exercise_id = ?", exerciseID).Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return mapper.SubmissionsFromGorm(submissions), nil
}

// GetSubmissionsByUser implements [SubmissionRepository].
func (s *submissionRepository) GetSubmissionsByUser(ctx context.Context, userID string) ([]models.Submission, error) {
	var submissions []lab_model.Submission
	err := s.db.WithContext(ctx).Model(&lab_model.Submission{}).Where("user_id = ?", userID).Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return mapper.SubmissionsFromGorm(submissions), nil
}

// UpdateSubmission implements [SubmissionRepository].
func (s *submissionRepository) UpdateSubmission(ctx context.Context, req *models.Submission) (*models.Submission, error) {
	submission := mapper.SubmissionToGorm(req)
	err := s.db.WithContext(ctx).Model(&lab_model.Submission{}).Where("id = ?", req.ID).Updates(submission).Error
	if err != nil {
		return nil, err
	}
	return mapper.SubmissionFromGorm(submission), nil
}
