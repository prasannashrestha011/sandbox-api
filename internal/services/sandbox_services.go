package services

import (
	"context"
	"log"

	"github.com/google/uuid"

	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/sandbox/core"
	services_validators "main/internal/services/validators"
	sandbox_type "main/internal/types"
)

// SandboxService exposes business operations for sandbox sessions.
type SandboxService interface {
	Create(ctx context.Context, sandbox *model.Sandbox) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Sandbox, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Sandbox, error)
	ListByUserID(ctx context.Context, userID string) ([]model.Sandbox, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status sandbox_type.SandboxState) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type sandboxService struct {
	repo          repository.SandboxRepository
	imageRepo     repository.DockerImageRepository
	sandboxClient core.SandboxClient
}

// NewSandboxService returns a service backed by a SandboxRepository.
func NewSandboxService(repo repository.SandboxRepository, imageRepo repository.DockerImageRepository, sandboxClient core.SandboxClient) SandboxService {
	return &sandboxService{repo: repo, imageRepo: imageRepo, sandboxClient: sandboxClient}
}

func (s *sandboxService) Create(ctx context.Context, sandbox *model.Sandbox) error {
	// Validate and cap the sandbox resource limits
	services_validators.ValidateAndCapSandboxLimits(sandbox)

	// call the sandbox client, the client returns the container id that to be stored in the database.
	//create method populates the model
	log.Println("Fetching image details for image ID: ", sandbox.ImageID)
	image, err := s.imageRepo.FindByID(ctx, sandbox.ImageID)
	if err != nil {
		return err
	}
	sandbox.Image = *image

	log.Println("Creating a sandbox")
	err = s.sandboxClient.Create(ctx, sandbox)
	if err != nil {
		log.Println("Sandbox client: ", err.Error())
		return err
	}
	err = s.repo.Create(ctx, sandbox)
	if err != nil {
		log.Println("Sandbox repository: ", err.Error())
		return err
	}

	return nil
}

func (s *sandboxService) GetByID(ctx context.Context, id uuid.UUID) (*model.Sandbox, error) {
	sandbox, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, postgres_error.MapError(err, "fetch sandbox by id", "sandbox")
	}
	return sandbox, nil
}

func (s *sandboxService) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Sandbox, error) {
	sandbox, err := s.repo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, postgres_error.MapError(err, "fetch sandbox by session ID", "sandbox")
	}
	return sandbox, nil
}

func (s *sandboxService) ListByUserID(ctx context.Context, userID string) ([]model.Sandbox, error) {
	sandboxes, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, postgres_error.MapError(err, "list sandboxes by user ID", "sandbox")
	}
	return sandboxes, nil
}

func (s *sandboxService) UpdateStatus(ctx context.Context, id uuid.UUID, status sandbox_type.SandboxState) error {
	return postgres_error.MapError(s.repo.UpdateStatus(ctx, id, status), "update sandbox status", "sandbox")
}

func (s *sandboxService) Delete(ctx context.Context, id uuid.UUID) error {
	return postgres_error.MapError(s.repo.Delete(ctx, id), "delete sandbox", "sandbox")
}
