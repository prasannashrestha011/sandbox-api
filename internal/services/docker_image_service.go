package services

import (
	"context"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/repository/model"
)

type DockerImageService interface {
	// CreateImage creates a new Docker image record in the database.
	CreateImage(imageTag string, environment string, userID string) error
	ListImages() ([]*model.DockerImage, error)
}

type dockerImageService struct {
	repo repository.DockerImageRepository
}

func NewDockerImageService(repo repository.DockerImageRepository) DockerImageService {
	return &dockerImageService{repo: repo}
}

func (s *dockerImageService) CreateImage(imageTag string, environment string, userID string) error {
	dockerImage := &model.DockerImage{
		ImageTag:    imageTag,
		Environment: environment,
		CreatedByID: userID,
	}
	err := postgres_error.MapError(s.repo.Create(context.Background(), dockerImage), "creating Docker image", "docker_image")
	if err != nil {
		return err
	}
	return nil
}

func (s *dockerImageService) ListImages() ([]*model.DockerImage, error) {
	list, err := s.repo.List(context.Background())
	if err != nil {
		return nil, postgres_error.MapError(err, "listing Docker images", "docker_image")
	}
	return list, nil
}
