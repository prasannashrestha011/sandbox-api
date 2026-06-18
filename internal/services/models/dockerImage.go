package models

import "time"

type DockerImage struct {
	ID          string
	ImageTag    string
	Runtime     string
	CreatedByID string
	CreatedBy   User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
