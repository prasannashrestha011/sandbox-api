package services

import (
	"context"
	"fmt"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/sandbox/core"
	"path"
	"strings"
)

type DockerImageService interface {
	// CreateImage creates a new Docker image record in the database.
	CreateImage(imageTag string, userID string) error
	ListImages() ([]*model.DockerImage, error)
}

type dockerImageService struct {
	repo          repository.DockerImageRepository
	sandboxClient core.SandboxClient
}

func NewDockerImageService(repo repository.DockerImageRepository, sandboxClient core.SandboxClient) DockerImageService {
	return &dockerImageService{repo: repo, sandboxClient: sandboxClient}
}

func langFromImageTag(imageTag string) (string, error) {
	// imageTag format: "language:version" or "registry/language:version"
	base := path.Base(imageTag) // strips registry prefix if any
	lang, _, found := strings.Cut(base, ":")
	if !found {
		return "", fmt.Errorf("invalid image tag format: %s", imageTag)
	}
	return lang, nil
}
func (s *dockerImageService) CreateImage(imageTag string, userID string) error {
	err := s.sandboxClient.PullImage(context.Background(), imageTag)
	if err != nil {
		return err
	}
	lang, err := langFromImageTag(imageTag)
	if err != nil {
		return fmt.Errorf("failed to extract language from image tag %s: %w", imageTag, err)
	}
	dockerImage := &model.DockerImage{
		ImageTag:    imageTag,
		CreatedByID: userID,
		Lang:        lang,
	}
	err = s.repo.Create(context.Background(), dockerImage)
	if err != nil {
		return postgres_error.MapError(err, "creating Docker image record", "docker_image")
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
