package lab_model

import (
	"main/internal/repository/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lab struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	Title       string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	Lang        string     `gorm:"not null"`
	Chapters    []Chapter  `gorm:"foreignKey:LabID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ContainerID string     `gorm:"not null"`
	CreatedByID string     `gorm:"not null;index;type:uuid"`
	CreatedBy   model.User `gorm:"foreignKey:CreatedByID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Tags       []Tag  `gorm:"many2many:lab_tags;"`
	Difficulty string `gorm:"not null"`
	IsPublic   bool   `gorm:"not null;default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (l *Lab) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}
