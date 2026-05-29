package repository

import (
	"context"
	"main/internal/repository/model"

	"gorm.io/gorm"
)

type DockerImageRepository interface {
	// Create inserts a new Docker image record into the database.
	Create(ctx context.Context, image *model.DockerImage) error
}

type dockerImageRepository struct {
	db *gorm.DB
}

// NewDockerImageRepository returns a GORM-backed DockerImageRepository.
func NewDockerImageRepository(db *gorm.DB) DockerImageRepository {
	return &dockerImageRepository{db: db}
}

func (r *dockerImageRepository) Create(ctx context.Context, image *model.DockerImage) error {
	return r.db.WithContext(ctx).Model(&model.DockerImage{}).Create(image).Error
}
