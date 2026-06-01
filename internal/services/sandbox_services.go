package services

import (
	"context"
	"log"

	"github.com/google/uuid"

	"main/internal/domain"
	"main/internal/dto"
	"main/internal/enums"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/sandbox/lang"
	service_mapper "main/internal/services/mapper"
	"main/internal/services/models"
	services_validators "main/internal/services/validators"
)

// SandboxService exposes business operations for sandbox sessions.
type SandboxService interface {
	Create(ctx context.Context, imageId string, sandbox *models.Sandbox) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Sandbox, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*models.Sandbox, error)
	ListByUserID(ctx context.Context, userID string) ([]models.Sandbox, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status enums.SandboxState) error

	//code exeuction
	ExecuteCode(ctx context.Context, containerID string, req *dto.ExecuteCodeRequest) (string, error)
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

func (s *sandboxService) Create(ctx context.Context, imageId string, sandbox *models.Sandbox) error {
	// Validate and cap the sandbox resource limits
	services_validators.ValidateAndCapSandboxLimits(sandbox)

	log.Println("Fetching image details for image ID: ", imageId)
	image, err := s.imageRepo.FindByID(ctx, imageId)
	if err != nil {
		return err
	}

	// Convert service model to repository model
	repoModel := service_mapper.ToRepoModel(sandbox)
	repoModel.ImageID = image.ID
	repoModel.Image = *image

	err = s.sandboxClient.Create(ctx, repoModel)
	if err != nil {
		log.Println("Sandbox client: ", err.Error())
		return err
	}

	err = s.repo.Create(ctx, repoModel)
	if err != nil {
		log.Println("Sandbox repository: ", err.Error())
		return err
	}

	// Sync back ID, SessionID, ContainerID, etc. created by repo/client
	sandbox.ID = repoModel.ID
	sandbox.SessionID = repoModel.SessionID
	sandbox.ContainerID = repoModel.ContainerID
	sandbox.ContainerName = repoModel.ContainerName
	// Note: We don't map Image over because service_models.DockerImage is not mapped yet

	return nil
}

func (s *sandboxService) GetByID(ctx context.Context, id uuid.UUID) (*models.Sandbox, error) {
	sandbox, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, postgres_error.MapError(err, "fetch sandbox by id", "sandbox")
	}
	return service_mapper.ToServiceModel(sandbox), nil
}

func (s *sandboxService) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*models.Sandbox, error) {
	sandbox, err := s.repo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, postgres_error.MapError(err, "fetch sandbox by session ID", "sandbox")
	}
	return service_mapper.ToServiceModel(sandbox), nil
}

func (s *sandboxService) ListByUserID(ctx context.Context, userID string) ([]models.Sandbox, error) {
	sandboxes, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, postgres_error.MapError(err, "list sandboxes by user ID", "sandbox")
	}

	var result []models.Sandbox
	for _, sb := range sandboxes {
		if mapped := service_mapper.ToServiceModel(&sb); mapped != nil {
			result = append(result, *mapped)
		}
	}
	return result, nil
}

func (s *sandboxService) UpdateStatus(ctx context.Context, id uuid.UUID, status enums.SandboxState) error {
	return postgres_error.MapError(s.repo.UpdateStatus(ctx, id, status), "update sandbox status", "sandbox")
}

func (s *sandboxService) ExecuteCode(ctx context.Context, containerID string, req *dto.ExecuteCodeRequest) (string, error) {

	cmd, err := lang.BuildCommand(req.Lang, req.Code)
	if err != nil {
		return "", domain.InvalidRequestError("unsupported language", nil)
	}

	return s.sandboxClient.ExecuteCode(ctx, containerID, cmd)
}
func (s *sandboxService) Delete(ctx context.Context, id uuid.UUID) error {
	return postgres_error.MapError(s.repo.Delete(ctx, id), "delete sandbox", "sandbox")
}
