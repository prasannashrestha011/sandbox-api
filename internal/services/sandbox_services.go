package services

import (
	"context"
	"log"

	"github.com/google/uuid"

	"main/internal/repository"
	"main/internal/repository/model"
	"main/internal/sandbox/core"
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
	return s.repo.FindByID(ctx, id)
}

func (s *sandboxService) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Sandbox, error) {
	return s.repo.FindBySessionID(ctx, sessionID)
}

func (s *sandboxService) ListByUserID(ctx context.Context, userID string) ([]model.Sandbox, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *sandboxService) UpdateStatus(ctx context.Context, id uuid.UUID, status sandbox_type.SandboxState) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *sandboxService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
