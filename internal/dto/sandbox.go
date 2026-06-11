package dto

import (
	"main/internal/enums"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateRequest payload for the sanbox creation
type CreateTemplateReq struct {
	// UserID string `json:"user_id,omitempty"`
	// Language Environment
	Lang string `json:"lang,omitempty"`
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
type UpdateRequest struct {
	MemoryLimit    *int64         `json:"memory_limit,omitempty"`
	CPULimit       *int64         `json:"cpu_limit,omitempty"`
	PidsLimit      *int64         `json:"pids_limit,omitempty"`
	SessionTimeout *time.Duration `json:"session_timeout,omitempty"`
	ExecTimeout    *time.Duration `json:"exec_timeout,omitempty"`
	NetworkMode    *string        `json:"network_mode,omitempty"`
}

// -- resource limit for sandbox --

func (r *CreateTemplateReq) Sanitize() {
	r.Lang = strings.TrimSpace(r.Lang)
	r.ImageID = strings.TrimSpace(r.ImageID)
	r.NetworkMode = strings.TrimSpace(r.NetworkMode)
}

func (r *CreateTemplateReq) Validate() error {
	r.Sanitize()
	var v ValidationErrors

	if r.Lang == "" && r.ImageID == "" {
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
	ImageTag    string `json:"image_tag,omitempty"`
	Environment string `json:"environment,omitempty"`
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

type SandboxTemplateResponse struct {
	ID             string        `json:"id"`
	UserID         string        `json:"user_id"`
	ImageTag       string        `json:"image_tag,omitempty"`
	Lang           string        `json:"lang,omitempty"`
	MemoryLimit    int64         `json:"memory_limit,omitempty"`
	CPULimit       int64         `json:"cpu_limit,omitempty"`
	PidsLimit      int64         `json:"pids_limit,omitempty"`
	SessionTimeout time.Duration `json:"session_timeout,omitempty"`
	ExecTimeout    time.Duration `json:"exec_timeout,omitempty"`
	NetworkMode    string        `json:"network_mode,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
}

type CreateSessionReq struct {
	TemplateID string `json:"template_id,omitempty"`
}

func (r *CreateSessionReq) Sanitize() {
	r.TemplateID = strings.TrimSpace(r.TemplateID)
}

func (r *CreateSessionReq) Validate() error {
	r.Sanitize()
	var v ValidationErrors

	if r.TemplateID == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "template_id",
			Message: "template_id is required",
		})
	}
	if _, err := uuid.Parse(r.TemplateID); err != nil {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "template_id",
			Message: "template_id must be a valid UUID",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

type SandboxSessionResponse struct {
	SessionID  string             `json:"session_id"`
	TemplateID string             `json:"template_id"`
	Status     enums.SandboxState `json:"status"`
	CreatedAt  time.Time          `json:"created_at"`
	ExpiresAt  time.Time          `json:"expires_at"`
}

type SandboxExecReq struct {
	SessionID string `json:"session_id,omitempty"`
	Command   string `json:"command,omitempty"`
}

type SandboxExecResponse struct {
	Stdout   string        `json:"stdout"`
	Stderr   string        `json:"stderr"`
	ExitCode int           `json:"exit_code"`
	ExecTime time.Duration `json:"exec_time"`
}
