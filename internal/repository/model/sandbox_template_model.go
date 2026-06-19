package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sandbox is the GORM model for a sandbox session.
type SandboxTemplate struct {
	ID      string      `gorm:"type:uuid;primaryKey"`
	UserID  string      `gorm:"not null;index"`
	Lang    string      `gorm:"not null"`
	ImageID string      `gorm:"type:uuid;not null"`
	Image   DockerImage `gorm:"foreignKey:ImageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	MemoryLimit    int64         `gorm:"not null"`
	CPULimit       int64         `gorm:"not null"`
	PidsLimit      int64         `gorm:"not null"`
	SessionTimeout time.Duration `gorm:"not null"`
	ExecTimeout    time.Duration `gorm:"not null"`
	NetworkMode    string        `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate ensures UUIDs are set before persisting.
func (s *SandboxTemplate) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
