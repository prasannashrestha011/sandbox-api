package lab_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Submission struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	ExerciseID  string `gorm:"not null;index;type:uuid"`
	UserID      string `gorm:"not null;index;type:uuid"`
	Code        string `gorm:"not null;type:text"`
	Output      string `gorm:"type:text"`
	Status      string `gorm:"not null"` // pending / passed / failed / error
	SubmittedAt time.Time
}

func (s *Submission) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
