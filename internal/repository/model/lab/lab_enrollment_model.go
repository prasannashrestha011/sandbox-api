package lab_model

import (
	"main/internal/repository/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabEnrollment struct {
	ID             string     `gorm:"primaryKey;type:uuid"`
	EnrolledUserID string     `gorm:"not null;type:uuid;uniqueIndex:uniq_user_lab"`
	User           model.User `gorm:"foreignKey:EnrolledUserID;references:UserID;constraint:OnDelete:CASCADE;"`
	LabID          string     `gorm:"not null;type:uuid;uniqueIndex:uniq_user_lab"`
	Lab            Lab        `gorm:"foreignKey:LabID;constraint:OnDelete:CASCADE;"`
	Status         string     `gorm:"not null;default:'enrolled'"` // enrolled | in_progress | completed
	ProgressPct    int        `gorm:"not null;default:0"`
	EnrolledAt     time.Time
	CompletedAt    *time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (l *LabEnrollment) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	return nil
}
