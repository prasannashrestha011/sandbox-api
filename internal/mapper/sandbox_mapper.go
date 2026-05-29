package mapper

import (
	"time"

	"main/internal/repository/model"
	sandbox_type "main/internal/types"
)

// SandboxCreateRequestToModel maps an API create request to a persistence model.
func SandboxCreateRequestToModel(req sandbox_type.CreateRequest, now time.Time) *model.Sandbox {
	expiresAt := now.Add(req.SessionTimeout)
	return &model.Sandbox{
		UserID:         req.UserID,
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
	}
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
