package model

import (
	"time"
)

type SandboxSession struct {
	ID string `gorm:"type:uuid;primaryKey"`

	// ownership
	UserID string `gorm:"not null;index"`

	// link to template
	TemplateID string `gorm:"not null;index"`

	Lang string `gorm:"not null"`
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
