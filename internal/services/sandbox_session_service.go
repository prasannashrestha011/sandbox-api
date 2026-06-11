package services

import (
	"context"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/enums"
	postgres_error "main/internal/infra/postgres"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/sandbox/lang"
	"main/internal/services/mapper"
)

type SandboxSessionService interface {
	// CreateSession creates a new sandbox session for a given user and template.
	CreateSession(ctx context.Context, templateID string) (*dto.SandboxSessionResponse, error)

	// GetActiveSession retrieves the active sandbox session for a given user and template.
	GetActiveSession(ctx context.Context, userID string, templateID string) (*dto.SandboxSessionResponse, error)

	// TerminateSession terminates an active sandbox session by its ID.
	TerminateSession(ctx context.Context, sessionID string) error
}

type sandboxSessionService struct {
	repo          repository.SandboxRepository
	template_repo repository.SandboxTemplateRepository
	client        core.SandboxClient
}

func NewSandboxSessionService(repo repository.SandboxRepository, client core.SandboxClient) SandboxSessionService {
	return &sandboxSessionService{repo: repo, client: client}
}

func (s *sandboxSessionService) CreateSession(ctx context.Context, templateID string) (*dto.SandboxSessionResponse, error) {
	session, err := mapper.ToSessionRequest(ctx, templateID)
	if err != nil {
		return nil, err
	}
	activeSession, err := s.repo.FindActiveSessionByUserAndTemplate(ctx, session.UserID, session.TemplateID)
	if err == nil && activeSession != nil {
		return mapper.ToSessionResponse(activeSession), nil
	}

	template, err := s.template_repo.FindByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	containerID, containerName, err := s.client.Create(ctx, template)
	if err != nil {
		return nil, err
	}

	session.SessionTimeout = template.SessionTimeout
	session.ExecTimeout = template.ExecTimeout
	session.Lang = template.Lang
	session.ContainerID = containerID
	session.ContainerName = containerName

	createdSession, err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, postgres_error.MapError(err, "create sandbox session", "sandbox_session")
	}
	return mapper.ToSessionResponse(createdSession), nil
}
func (s *sandboxSessionService) GetActiveSession(ctx context.Context, userID string, templateID string) (*dto.SandboxSessionResponse, error) {
	session, err := s.repo.FindActiveSessionByUserAndTemplate(ctx, userID, templateID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSessionResponse(session), nil
}
func (s *sandboxSessionService) ExecuteCommand(ctx context.Context, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error) {
	execModel, err := mapper.ToSandboxExecutionModel(ctx, req)
	if err != nil {
		return nil, err
	}
	session, err := s.repo.FindActiveSessionByUserAndTemplate(ctx, execModel.UserID, execModel.SessionID)

	if err != nil {
		return nil, err
	}
	if session.Status != enums.StateActive {
		return nil, domain.InvalidRequestError("lab session expired", nil)
	}

	cmd, err := lang.BuildCommand(session.Lang, execModel.Command)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.ExecuteCode(ctx, session.ContainerID, cmd)
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
