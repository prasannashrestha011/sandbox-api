package services

import (
	"context"
	"main/internal/repository"
	"main/internal/repository/model"
)

type DockerImageService interface {
	// CreateImage creates a new Docker image record in the database.
	CreateImage(imageTag string, userID string) error
}

type dockerImageService struct {
	repo repository.DockerImageRepository
}

func NewDockerImageService(repo repository.DockerImageRepository) DockerImageService {
	return &dockerImageService{repo: repo}
}

func (s *dockerImageService) CreateImage(imageTag string, userID string) error {
	dockerImage := &model.DockerImage{
		ImageTag:    imageTag,
		CreatedByID: userID,
	}
	return s.repo.Create(context.Background(), dockerImage)
}
