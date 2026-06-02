package services

import (
	"context"
	"net/http"

	"main/internal/domain"
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/services/mapper"
)

type LabService interface {
	CreateLab(ctx context.Context, req *dto.CreateLabRequest) (*dto.LabResponse, error)
	GetLabByID(ctx context.Context, id string) (*dto.LabResponse, error)
	DeleteLab(ctx context.Context, id string) error
}

type labService struct {
	repo repository.LabRepository
}

func NewLabService(repo repository.LabRepository) LabService {
	return &labService{
		repo: repo,
	}
}

func (s *labService) CreateLab(ctx context.Context, req *dto.CreateLabRequest) (*dto.LabResponse, error) {
	labModel := mapper.ToLabModel(req, ctx)

	ok, _ := s.repo.ValidateContainerID(ctx, labModel.ContainerID)
	if !ok {
		return nil, domain.NewAppError(http.StatusBadRequest, domain.CodeDockerImageNotFound, "specified docker image not found", nil, nil)
	}
	if err := s.repo.Create(ctx, labModel); err != nil {
		return nil, postgres_error.MapError(err, "CreateLab", "Lab")
	}

	return mapper.ToLabResponse(labModel), nil
}

func (s *labService) GetLabByID(ctx context.Context, id string) (*dto.LabResponse, error) {
	labModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, postgres_error.MapError(err, "GetLabByID", "Lab")
	}

	return mapper.ToLabResponse(labModel), nil
}

func (s *labService) DeleteLab(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return postgres_error.MapError(err, "DeleteLab", "Lab")
	}
	return nil
}
