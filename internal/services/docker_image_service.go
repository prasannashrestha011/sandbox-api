package services

import (
	"context"
	"fmt"
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/services/mapper"
)

type DockerImageService interface {
	// CreateImage creates a new Docker image record in the database.
	CreateImage(ctx context.Context, req *dto.CreateImageRequest) error
	ListImages(ctx context.Context) ([]*dto.DockerImageResponse, error)
}

type dockerImageService struct {
	repo          repository.DockerImageRepository
	sandboxClient core.SandboxClient
}

func NewDockerImageService(repo repository.DockerImageRepository, sandboxClient core.SandboxClient) DockerImageService {
	return &dockerImageService{repo: repo, sandboxClient: sandboxClient}
}

func (s *dockerImageService) CreateImage(ctx context.Context, req *dto.CreateImageRequest) error {
	dockerImage := mapper.ToDockerImageModel(ctx, req)
	err := s.sandboxClient.PullImage(context.Background(), dockerImage.ImageTag)
	if err != nil {
		return fmt.Errorf("failed to extract language from image tag %s: %w", req.ImageTag, err)
	}
	err = s.repo.Create(context.Background(), dockerImage)
	if err != nil {
		return postgres_error.MapError(err, "creating Docker image record", "docker_image")
	}
	return nil
}

func (s *dockerImageService) ListImages(ctx context.Context) ([]*dto.DockerImageResponse, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, postgres_error.MapError(err, "listing Docker images", "docker_image")
	}
	return mapper.ToDockerImageResponses(list), nil
}
