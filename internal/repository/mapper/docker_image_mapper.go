package mapper

import (
	"main/internal/repository/model"
	"main/internal/services/models"
)

func DockerImageToGom(d *models.DockerImage) *model.DockerImage {
	if d == nil {
		return nil
	}

	return &model.DockerImage{
		ID:          d.ID,
		ImageTag:    d.ImageTag,
		Lang:        d.Lang,
		CreatedByID: d.CreatedByID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func DockerImageFromGom(d *model.DockerImage) *models.DockerImage {
	if d == nil {
		return nil
	}

	return &models.DockerImage{
		ID:          d.ID,
		ImageTag:    d.ImageTag,
		Lang:        d.Lang,
		CreatedByID: d.CreatedByID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func DockerImagesFromGom(d []*model.DockerImage) []*models.DockerImage {
	if d == nil {
		return nil
	}

	images := make([]*models.DockerImage, len(d))
	for i, v := range d {
		images[i] = DockerImageFromGom(v)
	}
	return images
}
