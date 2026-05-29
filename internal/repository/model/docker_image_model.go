package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DockerImage struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	ImageTag    string `gorm:"unique;not null"`
	CreatedByID string `gorm:"not null;index;type:uuid"`
	CreatedBy   User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*DockerImage) TableName() string {
	return "docker_images"
}
func (d *DockerImage) BeforeCreate(tx *gorm.DB) (err error) {
	log.Println("generating the uuid")
	d.ID = uuid.New().String()
	return
}
