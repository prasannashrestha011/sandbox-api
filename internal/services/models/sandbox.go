package models

import (
	"main/internal/enums"
	"time"

	"github.com/google/uuid"
)

type Sandbox struct {
	ID             uuid.UUID
	UserID         string
	Environment    string
	Image          DockerImage
	MemoryLimit    int64
	CPULimit       int64
	PidsLimit      int64
	SessionTimeout time.Duration
	ExecTimeout    time.Duration
	NetworkMode    string

	ContainerName string
	ContainerID   string
	SessionID     uuid.UUID
	Status        enums.SandboxState
	ExpiresAt     time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
