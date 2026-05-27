package model

import (
	"main/internal/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID    uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Fullname  string         `gorm:"type:varchar(100);not null"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	Role      types.Role     `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserID = uuid.New()
	return nil
}
