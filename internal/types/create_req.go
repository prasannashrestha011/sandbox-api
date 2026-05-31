package types

import (
	"strings"
	"time"

	"main/internal/dto"
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

func (r *CreateRequest) Sanitize() {
	r.Environment = strings.TrimSpace(r.Environment)
	r.ImageID = strings.TrimSpace(r.ImageID)
	r.NetworkMode = strings.TrimSpace(r.NetworkMode)
}

func (r *CreateRequest) Validate() error {
	r.Sanitize()
	var v dto.ValidationErrors

	if r.Environment == "" && r.ImageID == "" {
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "environment",
			Message: "either environment or image_id is required",
		})
	}

	if r.MemoryLimit < 0 {
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "memory_limit",
			Message: "memory_limit cannot be negative",
		})
	}

	if r.CPULimit < 0 {
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "cpu_limit",
			Message: "cpu_limit cannot be negative",
		})
	}

	if r.PidsLimit < 0 {
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "pids_limit",
			Message: "pids_limit cannot be negative",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

type UpdateStatusRequest struct {
	Status SandboxState `json:"status"`
}

type CreateImageRequest struct {
	// Docker Image ID
	ImageTag string `json:"image_tag,omitempty"`
}

func (r *CreateImageRequest) Sanitize() {
	r.ImageTag = strings.TrimSpace(r.ImageTag)
}

func (r *CreateImageRequest) Validate() error {
	r.Sanitize()
	var v dto.ValidationErrors

	if r.ImageTag == "" {
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "image_tag",
			Message: "image_tag is required",
		})
	} else if !strings.Contains(r.ImageTag, ":") {
		// Assuming the required format is image:tag like ubuntu:general
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "image_tag",
			Message: "image_tag must be in a valid format containing a tag (e.g., name:general)",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}
