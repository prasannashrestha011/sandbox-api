package lab_services

import (
	"context"
	"main/internal/dto"
	repository "main/internal/repository/lab"
	"main/internal/services/mapper"
)

type SubmissionService interface {
	CreateSubmission(ctx context.Context, exerciseID string, req *dto.SubmissionRequest) (*dto.SubmissionResponse, error)
	GetSubmission(ctx context.Context, id string) (*dto.SubmissionResponse, error)
	ListSubmissions(ctx context.Context, exerciseID string) ([]dto.SubmissionResponse, error)
	UpdateSubmission(ctx context.Context, exerciseID string, req *dto.SubmissionRequest) (*dto.SubmissionResponse, error)
	DeleteSubmission(ctx context.Context, id string) error
}

type submissionService struct {
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
	}
}

func (s *submissionService) CreateSubmission(ctx context.Context, exerciseID string, req *dto.SubmissionRequest) (*dto.SubmissionResponse, error) {
	submissonModel := mapper.ToSubmissionModel(ctx, exerciseID, req)
	submission, err := s.submissionRepo.CreateSubmission(ctx, submissonModel)
	if err != nil {
		return nil, err
	}
	return mapper.ToSubmissionResponse(submission), nil
}

func (s *submissionService) DeleteSubmission(ctx context.Context, id string) error {
	err := s.submissionRepo.DeleteSubmission(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *submissionService) GetSubmission(ctx context.Context, id string) (*dto.SubmissionResponse, error) {
	submission, err := s.submissionRepo.GetSubmissionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToSubmissionResponse(submission), nil
}

func (s *submissionService) ListSubmissions(ctx context.Context, exerciseID string) ([]dto.SubmissionResponse, error) {
	submissions, err := s.submissionRepo.GetSubmissionsByExercise(ctx, exerciseID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSubmissionResponses(submissions), nil
}

func (s *submissionService) UpdateSubmission(ctx context.Context, exerciseID string, req *dto.SubmissionRequest) (*dto.SubmissionResponse, error) {
	submissionModel := mapper.ToSubmissionModel(ctx, exerciseID, req)
	submission, err := s.submissionRepo.UpdateSubmission(ctx, submissionModel)
	if err != nil {
		return nil, err
	}
	return mapper.ToSubmissionResponse(submission), nil
}
