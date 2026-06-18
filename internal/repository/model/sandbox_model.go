package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SandboxSession struct {
	ID string `gorm:"type:uuid;primaryKey"`

	// ownership
	UserID string `gorm:"not null;index"`

	// link to template
	TemplateID string `gorm:"not null;index"`

	Runtime string `gorm:"not null"`
	// runtime container info
	ContainerID   string
	ContainerName string

	// lifecycle
	Status string `gorm:"not null;index"` // running, stopped, expired, failed

	// time control
	StartedAt time.Time
	ExpiresAt time.Time
	EndedAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *SandboxSession) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	return nil
}
