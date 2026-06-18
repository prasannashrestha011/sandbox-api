package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DockerImage struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	ImageTag    string    `gorm:"unique;not null" json:"image_tag"`
	Runtime     string    `gorm:"not null" json:"runtime"`
	CreatedByID string    `gorm:"not null;index;type:uuid" json:"created_by_id"`
	CreatedBy   User      `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (*DockerImage) TableName() string {
	return "docker_images"
}
func (d *DockerImage) BeforeCreate(tx *gorm.DB) (err error) {
	log.Println("generating the uuid")
	d.ID = uuid.New().String()
	return
}
