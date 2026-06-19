package services

import (
	"context"
	"log"

	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/services/models"
	services_validators "main/internal/services/validators"
)

// SandboxTemplateService exposes business operations for sandbox templates.
type SandboxTemplateService interface {
	Create(ctx context.Context, imageId string, sandbox *models.SandboxTemplate) (*models.SandboxTemplate, error)
	GetByID(ctx context.Context, id string) (*models.SandboxTemplate, error)
	ListByUserID(ctx context.Context, userID string) ([]models.SandboxTemplate, error)
	UpdateDetails(ctx context.Context, id string, updates map[string]interface{}) error
}

type sandboxTemplateService struct {
	repo      repository.SandboxTemplateRepository
	imageRepo repository.DockerImageRepository
}

// NewSandboxTemplateService returns a service backed by a SandboxTemplateRepository.
func NewSandboxTemplateService(repo repository.SandboxTemplateRepository, imageRepo repository.DockerImageRepository, sandboxClient core.SandboxClient) SandboxTemplateService {
	return &sandboxTemplateService{repo: repo, imageRepo: imageRepo}
}

func (s *sandboxTemplateService) Create(ctx context.Context, imageId string, sandbox *models.SandboxTemplate) (*models.SandboxTemplate, error) {
	// Validate and cap the sandbox resource limits
	services_validators.ValidateAndCapSandboxLimits(sandbox)

	image, err := s.imageRepo.FindByID(ctx, imageId)
	if err != nil {
		return nil, err
	}
	sandbox.Image.ImageTag = image.ImageTag
	sandbox.Image.ID = image.ID
	sandbox.Image.Lang = image.Lang
	sandbox.Lang = image.Lang
	log.Println("Creating sandbox with image: ", image.ImageTag)
	createdSandbox, err := s.repo.Create(ctx, sandbox)
	if err != nil {
		log.Println("Sandbox repository: ", err.Error())
		return nil, err
	}
	*sandbox = *createdSandbox
	return sandbox, nil
}

func (s *sandboxTemplateService) GetByID(ctx context.Context, id string) (*models.SandboxTemplate, error) {
	sandbox, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, postgres_error.MapError(err, "fetch sandbox template by id", "sandbox_template")
	}
	return sandbox, nil
}

func (s *sandboxTemplateService) ListByUserID(ctx context.Context, userID string) ([]models.SandboxTemplate, error) {
	sandboxes, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, postgres_error.MapError(err, "list sandbox templates by user ID", "sandbox_template")
	}
	return sandboxes, nil
}

func (s *sandboxTemplateService) UpdateDetails(ctx context.Context, id string, updates map[string]interface{}) error {
	err := postgres_error.MapError(s.repo.UpdateDetails(ctx, id, updates), "update sandbox details", "sandbox")
	if err != nil {
		return err
	}
	return nil
}
