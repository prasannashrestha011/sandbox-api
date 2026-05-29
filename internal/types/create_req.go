package types

import (
	"time"
)

// CreateRequest payload for the sanbox creation
type CreateRequest struct {
	// UserID string `json:"user_id,omitempty"`
	// Language Environment
	Environment string `json:"environment,omitempty"`
	// Docker Image ID
	ImageID string `json:"image_id,omitempty"`

	// -- resource limit for sandbox --
	// Memory limit
	MemoryLimit int64 `json:"memory_limit,omitempty"`
	// CPU usage limit
	CPULimit int64 `json:"cpu_limit,omitempty"`
	// Number of processes limit
	PidsLimit int64 `json:"pids_limit,omitempty"`

	// Sandbox session time expiration
	SessionTimeout time.Duration `json:"session_timeout,omitempty"` // per-user session timeout
	// Executed Command duration
	ExecTimeout time.Duration `json:"exec_timeout,omitempty"` // per-command session timeout

	// Sandbox network mode
	NetworkMode string `json:"network_mode,omitempty"`
}

type UpdateStatusRequest struct {
	Status SandboxState `json:"status"`
}

type CreateImageRequest struct {
	// Docker Image ID
	ImageTag string `json:"image_tag,omitempty"`
}
