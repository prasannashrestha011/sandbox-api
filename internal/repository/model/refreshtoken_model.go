package model

import "time"

// RefreshToken maps to the `refresh_tokens` table.
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:text;not null;index"`
	TokenHash string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	Revoked   bool      `gorm:"default:false"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
