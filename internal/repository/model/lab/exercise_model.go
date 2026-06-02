package lab_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exercise struct {
	ID             string `gorm:"primaryKey;type:uuid"`
	LabID          string `gorm:"not null;index;type:uuid"`
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	StarterCode    string `gorm:"not null"`
	ExpectedOutput string `gorm:"not null"`
	Hints          string `gorm:"type:text"`
	OrderIndex     int    `gorm:"not null;default:0"`
	Solution       string `gorm:"type:text"`
	MaxAttempts    int    `gorm:"default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (e *Exercise) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}
