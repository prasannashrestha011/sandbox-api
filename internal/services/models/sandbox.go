package models

import (
	"main/internal/enums"
	"time"
)

// sandbox blueprint
type SandboxTemplate struct {
	ID             string
	UserID         string
	Image          DockerImage
	Lang           string
	MemoryLimit    int64
	CPULimit       int64
	PidsLimit      int64
	SessionTimeout time.Duration
	ExecTimeout    time.Duration
	NetworkMode    string
}
type SandboxInstance struct {
	ID         string
	PoolID     string
	TemplateID string
	UserID     string

	ContainerID   string
	ContainerName string

	Lang   string
	Status enums.SandboxState

	SessionTimeout time.Duration
	ExecTimeout    time.Duration
	LastUsedAt     time.Time
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
