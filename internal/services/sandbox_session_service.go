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
	"time"

	"github.com/hibiken/asynq"
)

type SandboxSessionService interface {
	// CreateSession creates a new sandbox session for a given user and template.
	CreateSession(ctx context.Context, templateID string) (*dto.SandboxSessionResponse, error)

	// GetActiveSession retrieves the active sandbox session for a given user and template.
	GetActiveSession(ctx context.Context, userID string, templateID string) (*dto.SandboxSessionResponse, error)

	// ExecuteCommand executes a command within an active sandbox session and returns the result.
	ExecuteCommand(ctx context.Context, sessionID string, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error)

	// TerminateSession terminates an active sandbox session by its ID.
	TerminateSession(ctx context.Context, sessionID string) error
}

type sandboxSessionService struct {
	repo          repository.SandboxRepository
	template_repo repository.SandboxTemplateRepository
	dockerclient  core.SandboxClient
	asynqclient   *asynq.Client
}

func NewSandboxSessionService(repo repository.SandboxRepository, templateRepo repository.SandboxTemplateRepository, dockerclient core.SandboxClient, asynqclient *asynq.Client) SandboxSessionService {
	return &sandboxSessionService{repo: repo, template_repo: templateRepo, dockerclient: dockerclient, asynqclient: asynqclient}
}

func (s *sandboxSessionService) CreateSession(ctx context.Context, templateID string) (*dto.SandboxSessionResponse, error) {
	session, err := mapper.ToSessionRequest(ctx, templateID)
	if err != nil {
		return nil, err
	}
	activeSession, err := s.repo.FindActiveSessionByUser(ctx, session.UserID, session.TemplateID)
	if err == nil && activeSession != nil {
		return mapper.ToSessionResponse(activeSession), nil
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

	session.SessionTimeout = template.SessionTimeout
	session.ExecTimeout = template.ExecTimeout
	session.Runtime = template.Runtime
	session.ContainerID = containerID
	session.ContainerName = containerName
	session.Status = enums.StateActive
	session.ExpiresAt = time.Now().Add(template.SessionTimeout * 3600 * time.Second)

	// payload := &dto.SandboxCleanupPayload{
	// 	ContainerID: containerID,
	// 	SessionID:   session.ID,
	// }

	// task, err := jobs.NewSandboxCleanupTask(payload)
	// if err != nil {
	// 	log.Printf("Error creating sandbox cleanup task for session %s: %v", session.ID, err)
	// 	return nil, err
	// }
	// if _, err = s.asynqclient.Enqueue(task, asynq.ProcessIn(30*time.Second)); err != nil {
	// 	log.Printf("Error enqueuing sandbox cleanup task for session %s: %v", session.ID, err)
	// 	return nil, err

	// }

	createdSession, err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, postgres_error.MapError(err, "create sandbox session", "sandbox_session")
	}
	return mapper.ToSessionResponse(createdSession), nil
}
func (s *sandboxSessionService) GetActiveSession(ctx context.Context, userID string, templateID string) (*dto.SandboxSessionResponse, error) {
	session, err := s.repo.FindActiveSessionByUser(ctx, userID, templateID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSessionResponse(session), nil
}
func (s *sandboxSessionService) ExecuteCommand(ctx context.Context, sessionID string, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error) {
	execModel, err := mapper.ToSandboxExecutionModel(ctx, sessionID, req)
	if err != nil {
		return nil, err
	}
	session, err := s.repo.FindActiveSessionByUser(ctx, execModel.UserID, sessionID)
	if err != nil {
		return nil, err

	}
	if session.Status != enums.StateActive {
		return nil, domain.InvalidRequestError("lab session expired", nil)
	}

	cmd, err := lang.BuildCommand(session.Runtime, execModel.Command)
	if err != nil {
		return nil, err
	}
	resp, err := s.dockerclient.ExecuteCode(ctx, session.ContainerID, cmd)
	if err != nil {
		return nil, err
	}
	execModel.ExitCode = resp.ExitCode
	execModel.Stdout = resp.Stdout
	execModel.Stderr = resp.Stderr
	execModel.ExecTime = resp.ExecTime
	return resp, nil
}

func (s *sandboxSessionService) TerminateSession(ctx context.Context, sessionID string) error {
	return s.repo.UpdateStatus(ctx, sessionID, "terminated")
}
