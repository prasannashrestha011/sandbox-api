package mapper

import (
	"main/internal/repository/model"
	service_models "main/internal/services/models"
)

// ToRepoModel maps a service Sandbox model to a persistence repo Sandbox model.
func ToRepoModel(s *service_models.Sandbox) *model.Sandbox {
	if s == nil {
		return nil
	}
	return &model.Sandbox{
		ID:             s.ID,
		UserID:         s.UserID,
		Environment:    s.Environment,
		MemoryLimit:    s.MemoryLimit,
		CPULimit:       s.CPULimit,
		PidsLimit:      s.PidsLimit,
		SessionTimeout: s.SessionTimeout,
		ExecTimeout:    s.ExecTimeout,
		NetworkMode:    s.NetworkMode,
		ContainerName:  s.ContainerName,
		ContainerID:    s.ContainerID,
		SessionID:      s.SessionID,
		Status:         s.Status,
		ExpiresAt:      s.ExpiresAt,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		// Note: Image mapping can be added here if needed recursively.
		// You may just pass ImageID depending on how the repository handles it.
	}
}

// ToServiceModel maps a persistence repo Sandbox model to a service Sandbox model.
func ToServiceModel(r *model.Sandbox) *service_models.Sandbox {
	if r == nil {
		return nil
	}
	return &service_models.Sandbox{
		ID:             r.ID,
		UserID:         r.UserID,
		Environment:    r.Environment,
		MemoryLimit:    r.MemoryLimit,
		CPULimit:       r.CPULimit,
		PidsLimit:      r.PidsLimit,
		SessionTimeout: r.SessionTimeout,
		ExecTimeout:    r.ExecTimeout,
		NetworkMode:    r.NetworkMode,
		ContainerName:  r.ContainerName,
		ContainerID:    r.ContainerID,
		SessionID:      r.SessionID,
		Status:         r.Status,
		ExpiresAt:      r.ExpiresAt,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}
