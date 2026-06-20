package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WarmPool struct {
	ID string `gorm:"type:uuid;primaryKey"`

	TemplateID string `gorm:"not null;index"`

	MaxActive int `gorm:"not null"`

	Status string `gorm:"not null;index"` // active, inactive

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *WarmPool) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}

	return nil
}
