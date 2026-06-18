package mapper

import (
	"context"
	"net/http"
	"time"

	request_context "main/internal/context"
	"main/internal/dto"
	"main/internal/services/models"
)

// SandboxCreateRequestToServiceModel maps an API create request to a service model.
func ToSandboxTemplate(req dto.CreateTemplateReq, ctx context.Context, now time.Time) (*models.SandboxTemplate, error) {
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	return &models.SandboxTemplate{
		UserID:         userID.String(),
		Runtime:        req.Runtime,
		MemoryLimit:    req.MemoryLimit,
		CPULimit:       req.CPULimit,
		PidsLimit:      req.PidsLimit,
		SessionTimeout: req.SessionTimeout,
		ExecTimeout:    req.ExecTimeout,
		NetworkMode:    req.NetworkMode,
	}, nil
}

func ToSandboxTemplateResponse(sandbox *models.SandboxTemplate) *dto.SandboxTemplateResponse {
	if sandbox == nil {
		return &dto.SandboxTemplateResponse{}
	}

	return &dto.SandboxTemplateResponse{
		ID:             sandbox.ID,
		UserID:         sandbox.UserID,
		MemoryLimit:    sandbox.MemoryLimit,
		Runtime:        sandbox.Runtime,
		CPULimit:       sandbox.CPULimit,
		PidsLimit:      sandbox.PidsLimit,
		SessionTimeout: sandbox.SessionTimeout,
		ExecTimeout:    sandbox.ExecTimeout,
		NetworkMode:    sandbox.NetworkMode,
	}
}

func ToSandboxTemplateResponseList(sandboxes []models.SandboxTemplate) []*dto.SandboxTemplateResponse {
	responses := make([]*dto.SandboxTemplateResponse, len(sandboxes))
	for i, sb := range sandboxes {
		responses[i] = ToSandboxTemplateResponse(&sb)
	}
	return responses
}

func ToSessionRequest(ctx context.Context, templateID string) (*models.SandboxSession, error) {
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	return &models.SandboxSession{
		UserID:     userID.String(),
		TemplateID: templateID,
	}, nil
}
func ToSessionResponse(session *models.SandboxSession) *dto.SandboxSessionResponse {
	if session == nil {
		return &dto.SandboxSessionResponse{}
	}

	return &dto.SandboxSessionResponse{
		SessionID:  session.ID,
		TemplateID: session.TemplateID,
		Status:     session.Status,
		CreatedAt:  session.CreatedAt,
		ExpiresAt:  session.ExpiresAt,
	}
}

func ToSandboxExecutionModel(ctx context.Context, sessionID string, req *dto.SandboxExecReq) (*models.SandboxExecution, error) {

	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	return &models.SandboxExecution{
		UserID:  userID.String(),
		Command: req.Command,
	}, nil
}
func ToSandboxExecutionResponse(exec *models.SandboxExecution) *dto.SandboxExecResponse {
	if exec == nil {
		return &dto.SandboxExecResponse{}
	}

	return &dto.SandboxExecResponse{
		Stdout:   exec.Stdout,
		Stderr:   exec.Stderr,
		ExitCode: exec.ExitCode,
		ExecTime: exec.ExecTime,
	}
}

// // SandboxServiceModelToCreateResponse maps a service sandbox model to the create response payload.
// func SandboxServiceModelToCreateResponse(sandbox *models.SandboxTemplate) dto.CreateResponse {
// 	if sandbox == nil {
// 		return dto.CreateResponse{}
// 	}

// 	return dto.CreateResponse{
// 		ContainerID: sandbox.ContainerID,
// 		SessionID:   sandbox.SessionID,
// 		Status:      sandbox.Status,
// 		CreatedAt:   sandbox.CreatedAt,
// 		ExpiresAt:   sandbox.ExpiresAt,
// 		Error:       nil,
// 	}
// }
