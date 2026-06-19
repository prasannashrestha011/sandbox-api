package services

import (
	"context"
	"main/internal/domain"
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/sandbox/lang"
	"main/internal/services/mapper"
)

type DockerImageService interface {
	// CreateImage creates a new Docker image record in the database.
	CreateImage(ctx context.Context, req *dto.CreateImageRequest) (*dto.DockerImageResponse, error)
	ListImages(ctx context.Context) ([]*dto.DockerImageResponse, error)
}

type dockerImageService struct {
	repo          repository.DockerImageRepository
	sandboxClient core.SandboxClient
}

func NewDockerImageService(repo repository.DockerImageRepository, sandboxClient core.SandboxClient) DockerImageService {
	return &dockerImageService{repo: repo, sandboxClient: sandboxClient}
}

func (s *dockerImageService) CreateImage(ctx context.Context, req *dto.CreateImageRequest) (*dto.DockerImageResponse, error) {
	dockerImage := mapper.ToDockerImageModel(ctx, req)
	// after successfull validation, finally pull the  image
	cmd, err := lang.BuildCommand(req.Lang, "")
	if len(cmd) == 0 && err != nil {
		return nil, err
	}
	err = s.sandboxClient.PullImage(context.Background(), dockerImage.ImageTag)
	if err != nil {
		return nil, domain.InvalidRequestError("failed to pull docker image", nil)
	}
	createdImage, err := s.repo.Create(context.Background(), dockerImage)
	if err != nil {
		return nil, postgres_error.MapError(err, "creating Docker image record", "docker_image")
	}
	return mapper.ToDockerImageResponse(createdImage), nil
}

func (s *dockerImageService) ListImages(ctx context.Context) ([]*dto.DockerImageResponse, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, postgres_error.MapError(err, "listing Docker images", "docker_image")
	}
	return mapper.ToDockerImageResponses(list), nil
}
