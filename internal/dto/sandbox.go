package dto

import (
	"main/internal/enums"
	"strings"
	"time"

	"github.com/google/uuid"
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
	var v ValidationErrors

	if r.Environment == "" && r.ImageID == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "environment",
			Message: "either environment or image_id is required",
		})
	}

	if r.MemoryLimit < 0 {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "memory_limit",
			Message: "memory_limit cannot be negative",
		})
	}

	if r.CPULimit < 0 {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "cpu_limit",
			Message: "cpu_limit cannot be negative",
		})
	}

	if r.PidsLimit < 0 {
		v.Violations = append(v.Violations, FieldViolation{
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
	Status enums.SandboxState `json:"status"`
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
	var v ValidationErrors

	if r.ImageTag == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "image_tag",
			Message: "image_tag is required",
		})
	} else if !strings.Contains(r.ImageTag, ":") {
		// Assuming the required format is image:tag like ubuntu:general
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "image_tag",
			Message: "image_tag must be in a valid format containing a tag (e.g., name:general)",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

type ExecuteCodeRequest struct {
	Lang string `json:"lang"`
	Code string `json:"code"`
}

func (r *ExecuteCodeRequest) Sanitize() {
	r.Code = strings.TrimSpace(r.Code)
}

func (r *ExecuteCodeRequest) Validate() error {
	r.Sanitize()
	var v ValidationErrors

	if r.Code == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "code",
			Message: "code is required",
		})
	}
	if r.Lang == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "lang",
			Message: "lang is required",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

// CreateResponse : sandbox creation result
type CreateResponse struct {
	ContainerID string             `json:"container_id"`
	SessionID   uuid.UUID          `json:"session_id"`
	Status      enums.SandboxState `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	ExpiresAt   time.Time          `json:"expires_at"`
	Error       *string            `json:"error"`
}

// CodeResponse represents the response structure for code execution results.
type CodeResponse struct {
	CorrectResult  string `json:"correctResult"`
	ExecutedResult string `json:"executedResult"`
	Success        bool   `json:"success"`
}
