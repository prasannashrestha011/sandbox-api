package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DockerImage struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	ImageTag    string    `gorm:"unique;not null"`
	CreatedByID uuid.UUID `gorm:"not null;index"`
	CreatedBy   User      `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*DockerImage) TableName() string {
	return "docker_images"
}
func (d *DockerImage) BeforeCreate(tx *gorm.DB) (err error) {
	log.Println("generating the uuid")
	d.ID = uuid.New()
	return
}
