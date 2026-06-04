package lab_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chapter struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	LabID       string     `gorm:"not null;index;type:uuid"`
	Title       string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	OrderIndex  int        `gorm:"not null;default:0"`
	Exercises   []Exercise `gorm:"foreignKey:ChapterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c *Chapter) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}
