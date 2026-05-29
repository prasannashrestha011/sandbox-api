package services

import (
	"context"
	"main/internal/repository"
	"main/internal/repository/model"

	"github.com/google/uuid"
)

type DockerImageService interface {
	// CreateImage creates a new Docker image record in the database.
	CreateImage(imageID string, userID uuid.UUID) error
}

type dockerImageService struct {
	repo repository.DockerImageRepository
}

func NewDockerImageService(repo repository.DockerImageRepository) DockerImageService {
	return &dockerImageService{repo: repo}
}

func (s *dockerImageService) CreateImage(imageID string, userID uuid.UUID) error {
	dockerImage := &model.DockerImage{
		ImageTag:    imageID,
		CreatedByID: userID,
	}
	return s.repo.Create(context.Background(), dockerImage)
}
