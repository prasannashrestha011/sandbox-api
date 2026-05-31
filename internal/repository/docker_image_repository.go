package repository

import (
	"context"
	"main/internal/repository/model"

	"gorm.io/gorm"
)

type DockerImageRepository interface {
	// Create inserts a new Docker image record into the database.
	Create(ctx context.Context, image *model.DockerImage) error
	FindByID(ctx context.Context, id string) (*model.DockerImage, error)
	List(ctx context.Context) ([]*model.DockerImage, error)
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

func (r *dockerImageRepository) FindByID(ctx context.Context, id string) (*model.DockerImage, error) {
	var image model.DockerImage
	err := r.db.WithContext(ctx).Model(&model.DockerImage{}).Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *dockerImageRepository) List(ctx context.Context) ([]*model.DockerImage, error) {
	var images []*model.DockerImage
	err := r.db.WithContext(ctx).Model(&model.DockerImage{}).Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}
