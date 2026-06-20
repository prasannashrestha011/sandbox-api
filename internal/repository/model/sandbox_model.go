package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SandboxInstance struct {
	ID string `gorm:"type:uuid;primaryKey"`

	// ownership
	UserID string `gorm:"not null;index"`

	// link to template
	TemplateID string `gorm:"not null;index"`
	PoolID     string `gorm:"not null;index"`

	Lang string `gorm:"not null"`
	// runtime container info
	ContainerID string

	// lifecycle
	Status string `gorm:"not null;index"` // running, stopped, expired, failed

	LastUsed time.Time `gorm:"not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *SandboxInstance) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	return nil
}
