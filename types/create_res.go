package sandbox_request

import (
	"time"

	"github.com/google/uuid"
)

type SandBoxState string

const (
	StateActive   SandBoxState = "active"
	StateInActive SandBoxState = "inactive"
)

// CreateResponse : sandbox creation result
type CreateResponse struct {
	ContainerID string
	SessionID   uuid.UUID
	Status      SandBoxState
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Error       *string
}
