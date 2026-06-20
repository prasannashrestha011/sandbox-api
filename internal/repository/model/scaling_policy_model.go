package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScalingPolicy struct {
	ID string `gorm:"type:uuid;primaryKey"`

	WarmPoolID string `gorm:"not null;index"`

	// core thresholds
	MinIdleThreshold int `gorm:"not null"` // trigger scale up
	MaxIdleThreshold int `gorm:"not null"` // trigger scale down

	// behavior tuning
	ScaleUpStep   int `gorm:"not null"` // how many to create at once
	ScaleDownStep int `gorm:"not null"` // how many to remove at once

	CooldownSec int `gorm:"not null"` // prevents rapid scaling

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *ScalingPolicy) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
