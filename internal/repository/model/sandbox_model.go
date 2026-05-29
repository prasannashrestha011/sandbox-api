package model

import (
	sandbox_type "main/internal/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sandbox is the GORM model for a sandbox session.
type Sandbox struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey"`
	UserID         string        `gorm:"not null;index"`
	Environment    string        `gorm:"not null"`
	ImageID        string        `gorm:"not null"`
	MemoryLimit    int64         `gorm:"not null"`
	CPULimit       int64         `gorm:"not null"`
	PidsLimit      int64         `gorm:"not null"`
	SessionTimeout time.Duration `gorm:"not null"`
	ExecTimeout    time.Duration `gorm:"not null"`
	NetworkMode    string        `gorm:"not null"`

	ContainerID string                    `gorm:"index"`
	SessionID   uuid.UUID                 `gorm:"type:uuid;not null;uniqueIndex"`
	Status      sandbox_type.SandboxState `gorm:"type:varchar(16);not null"`
	ExpiresAt   time.Time                 `gorm:"not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate ensures UUIDs are set before persisting.
func (s *Sandbox) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.SessionID == uuid.Nil {
		s.SessionID = uuid.New()
	}
	return nil
}
