package lab_model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID   string `gorm:"primaryKey;type:uuid"`
	Name string `gorm:"not null;uniqueIndex"`
	Labs []Lab  `gorm:"many2many:lab_tags;"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
