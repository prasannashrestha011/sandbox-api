package mapper

import (
	"context"
	"net/http"
	"time"

	request_context "main/internal/context"
	"main/internal/dto"
	"main/internal/enums"
	"main/internal/services/models"
)

// SandboxCreateRequestToServiceModel maps an API create request to a service model.
func SandboxCreateRequestToServiceModel(req dto.CreateRequest, ctx context.Context, now time.Time) (*models.Sandbox, error) {
	expiresAt := now.Add(req.SessionTimeout)
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	return &models.Sandbox{
		UserID:         userID.String(),
		Environment:    req.Environment,
		MemoryLimit:    req.MemoryLimit,
		CPULimit:       req.CPULimit,
		PidsLimit:      req.PidsLimit,
		SessionTimeout: req.SessionTimeout,
		ExecTimeout:    req.ExecTimeout,
		NetworkMode:    req.NetworkMode,
		Status:         enums.StateActive,
		ExpiresAt:      expiresAt,
	}, nil
}

// SandboxServiceModelToCreateResponse maps a service sandbox model to the create response payload.
func SandboxServiceModelToCreateResponse(sandbox *models.Sandbox) dto.CreateResponse {
	if sandbox == nil {
		return dto.CreateResponse{}
	}

	return dto.CreateResponse{
		ContainerID: sandbox.ContainerID,
		SessionID:   sandbox.SessionID,
		Status:      sandbox.Status,
		CreatedAt:   sandbox.CreatedAt,
		ExpiresAt:   sandbox.ExpiresAt,
		Error:       nil,
	}
}
