package services

import (
	"context"
	"log"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/enums"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/sandbox/lang"
	"main/internal/services/mapper"
	"main/internal/services/models"

	"github.com/hibiken/asynq"
)

type SandboxInstanceService interface {
	// CreatePool creates a new sandbox pool for a given template.
	CreateInstance(ctx context.Context, templateID string) (*models.SandboxInstance, error)

	AcquireInstance(ctx context.Context, lang string) (*models.SandboxInstance, error)

	// ExecuteCommand executes a command within an active sandbox session and returns the result.
	ExecuteCommand(ctx context.Context, containerID string, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error)

	// TerminateSession terminates an active sandbox session by its ID.
	TerminateSession(ctx context.Context, sessionID string) error
}

type sandboxInstanceService struct {
	repo          repository.SandboxInstanceRepository
	template_repo repository.SandboxTemplateRepository
	dockerclient  core.SandboxClient
	asynqclient   *asynq.Client
}

func NewSandboxInstanceService(repo repository.SandboxInstanceRepository, templateRepo repository.SandboxTemplateRepository, dockerclient core.SandboxClient, asynqclient *asynq.Client) SandboxInstanceService {
	return &sandboxInstanceService{repo: repo, template_repo: templateRepo, dockerclient: dockerclient, asynqclient: asynqclient}
}

func (s *sandboxInstanceService) CreateInstance(ctx context.Context, templateID string) (*models.SandboxInstance, error) {
	session, err := mapper.ToSessionRequest(ctx, templateID)
	if err != nil {
		return nil, err
	}

	template, err := s.template_repo.FindByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	log.Printf("Creating sandbox session for user %s with template %s", session.UserID, template.Image.ImageTag)
	containerID, containerName, err := s.dockerclient.Create(ctx, template)
	if err != nil {
		return nil, err
	}

	session.ExecTimeout = template.ExecTimeout
	session.Lang = template.Lang
	session.ContainerID = containerID
	session.ContainerName = containerName
	session.Status = enums.StateActive
	instance, err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, postgres_error.MapError(err, "create sandbox session", "sandbox_session")
	}
	return instance, nil
}
func (s *sandboxInstanceService) AcquireInstance(ctx context.Context, lang string) (*models.SandboxInstance, error) {
	instance, err := s.repo.Acquire(ctx, lang)
	if err != nil {
		return nil, postgres_error.MapError(err, "acquire sandbox instance", "sandbox_instance")
	}
	return instance, nil
}
func (s *sandboxInstanceService) ExecuteCommand(ctx context.Context, containerID string, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error) {
	execModel, err := mapper.ToSandboxExecutionModel(ctx, containerID, req)
	if err != nil {
		return nil, err
	}

	cmd, err := lang.BuildCommand(execModel.Lang, execModel.Command)
	if err != nil {
		return nil, err
	}
	resp, err := s.dockerclient.ExecuteCode(ctx, execModel.ID, cmd)
	if err != nil {
		log.Println("error in sever")
		return nil, domain.InvalidRequestError("execution timeout", nil)
	}
	execModel.ExitCode = resp.ExitCode
	execModel.Stdout = resp.Stdout
	execModel.Stderr = resp.Stderr
	execModel.ExecTime = resp.ExecTime
	return resp, nil
}

func (s *sandboxInstanceService) TerminateSession(ctx context.Context, sessionID string) error {
	return s.repo.UpdateStatus(ctx, sessionID, "terminated")
}
