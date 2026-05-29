package types

import (
	"time"

	"github.com/google/uuid"
)

// CreateResponse : sandbox creation result
type CreateResponse struct {
	ContainerID string
	SessionID   uuid.UUID
	Status      SandboxState
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Error       *string
}
