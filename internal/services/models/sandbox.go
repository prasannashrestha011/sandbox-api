package models

import (
	"main/internal/enums"
	"time"
)

// sandbox blueprint
type SandboxTemplate struct {
	ID             string
	UserID         string
	Runtime        string
	Image          DockerImage
	MemoryLimit    int64
	CPULimit       int64
	PidsLimit      int64
	SessionTimeout time.Duration
	ExecTimeout    time.Duration
	NetworkMode    string
}
type SandboxSession struct {
	ID         string
	TemplateID string
	UserID     string

	ContainerID   string
	ContainerName string

	Runtime string
	Status  enums.SandboxState

	SessionTimeout time.Duration
	ExecTimeout    time.Duration
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SandboxExecution struct {
	ID        string
	UserID    string
	SessionID string

	Command   string
	Stdout    string
	Stderr    string
	ExitCode  int
	ExecTime  time.Duration
	CreatedAt time.Time
}
