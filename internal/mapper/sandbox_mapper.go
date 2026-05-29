package mapper

import (
	"context"
	"net/http"
	"time"

	request_context "main/internal/context"
	"main/internal/repository/model"
	sandbox_type "main/internal/types"
)

// SandboxCreateRequestToModel maps an API create request to a persistence model.
func SandboxCreateRequestToModel(req sandbox_type.CreateRequest, ctx context.Context, now time.Time) (*model.Sandbox, error) {
	expiresAt := now.Add(req.SessionTimeout)
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	return &model.Sandbox{
		UserID:         userID.String(),
		Environment:    req.Environment,
		ImageID:        req.ImageID,
		MemoryLimit:    req.MemoryLimit,
		CPULimit:       req.CPULimit,
		PidsLimit:      req.PidsLimit,
		SessionTimeout: req.SessionTimeout,
		ExecTimeout:    req.ExecTimeout,
		NetworkMode:    req.NetworkMode,
		Status:         sandbox_type.StateActive,
		ExpiresAt:      expiresAt,
	}, nil
}

// SandboxModelToCreateResponse maps a sandbox model to the create response payload.
func SandboxModelToCreateResponse(sandbox *model.Sandbox) sandbox_type.CreateResponse {
	if sandbox == nil {
		return sandbox_type.CreateResponse{}
	}

	return sandbox_type.CreateResponse{
		ContainerID: sandbox.ContainerID,
		SessionID:   sandbox.SessionID,
		Status:      sandbox.Status,
		CreatedAt:   sandbox.CreatedAt,
		ExpiresAt:   sandbox.ExpiresAt,
		Error:       nil,
	}
}
