package lab_model

import (
	"main/internal/repository/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabEnrollment struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	UserID      string     `gorm:"not null;index;type:uuid"`
	LabID       string     `gorm:"not null;index;type:uuid"`
	User        model.User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;"`
	Lab         Lab        `gorm:"foreignKey:LabID;constraint:OnDelete:CASCADE;"`
	Status      string     `gorm:"not null;default:'enrolled'"` // enrolled | in_progress | completed
	ProgressPct int        `gorm:"not null;default:0"`
	EnrolledAt  time.Time
	CompletedAt *time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (l *LabEnrollment) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	return nil
}
