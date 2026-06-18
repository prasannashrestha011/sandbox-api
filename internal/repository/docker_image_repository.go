package repository

import (
	"context"
	"main/internal/repository/mapper"
	gormodel "main/internal/repository/model"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type dockerImageRepository struct {
	db *gorm.DB
}

type DockerImageRepository interface {
	// Create inserts a new Docker image record into the database.
	Create(ctx context.Context, req *models.DockerImage) error
	FindByID(ctx context.Context, id string) (*models.DockerImage, error)
	List(ctx context.Context) ([]*models.DockerImage, error)
}

// NewDockerImageRepository returns a GORM-backed DockerImageRepository.
func NewDockerImageRepository(db *gorm.DB) DockerImageRepository {
	return &dockerImageRepository{db: db}
}

func (r *dockerImageRepository) Create(ctx context.Context, req *models.DockerImage) error {
	image := mapper.DockerImageToGom(req)
	return r.db.WithContext(ctx).Model(&gormodel.DockerImage{}).Create(image).Error
}

func (r *dockerImageRepository) FindByID(ctx context.Context, id string) (*models.DockerImage, error) {
	var image gormodel.DockerImage
	err := r.db.WithContext(ctx).Model(&gormodel.DockerImage{}).Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return mapper.DockerImageFromGom(&image), nil
}

func (r *dockerImageRepository) List(ctx context.Context) ([]*models.DockerImage, error) {
	var images []*gormodel.DockerImage
	err := r.db.WithContext(ctx).Model(&gormodel.DockerImage{}).Find(&images).Error
	if err != nil {
		return nil, err
	}
	return mapper.DockerImagesFromGom(images), nil
}
